log_format  detailed '$remote_addr - $remote_user [$time_local] "$request" '
                '$status $body_bytes_sent "$http_referer" '
                '"$http_user_agent" "$request_time" "$upstream_connect_time" "$upstream_response_time"';

proxy_headers_hash_bucket_size 128;
proxy_headers_hash_max_size 1024;

geo $mywhitelist {
  default 1;
  10.0.0.0/16 0;
  172.21.0.0/16 0;
#  182.156.218.114/32 0; ## India Office

  13.126.0.0/16 0;  ## Load Test IP
  13.251.0.0/16 0; ## load test IP

  172.31.0.0/16 0;      ## internal-office

#  192.168.0.0/16 0;     ## azure

}

map $mywhitelist $iplimit {
  1  $binary_remote_addr;
  0  "";
}

map $status $logger {
      ~^[23]  0;
      default 1;
}

limit_req_zone $iplimit  zone=fiveip:25m rate=5r/s;
limit_req_status 429;

server {
   listen 81;
   server_name boilerplate.your_domain.com;

   location / {
      return 301  https://boilerplate.your_domain.com$request_uri;
  }
}

server {
  listen 80;
  server_name  boilerplate.your_domain.com  boilerplate.service.azure;

  set $boilderplatenode http://127.0.0.1:9000;

  set $ffnode http://internal-service.service.aliyun.consul;
  root /var/www/boilerplate/public/tmpl/;

  proxy_set_header    X-Forwarded-For             $remote_addr;
  proxy_set_header    Host                        $host;
  proxy_set_header    X-Forwarded-Host            $host:$server_port;
  proxy_set_header    X-Forwarded-Server          $server_name;
  proxy_set_header    X-Forwarded-For             $remote_addr;
  proxy_set_header    X-Forwarded-Request-Uri     $request_uri;


  access_log  /var/log/nginx/scalyr.boilerplate.access.log main if=$logger;
  access_log /var/log/nginx/boilerplate.access.log detailed;
  error_log  /var/log/nginx/boilerplate.error.log;

  location ~ ^/v1/api/admin {
    limit_req zone=fiveip burst=10 nodelay;
    proxy_pass   $boilderplatenode;
  }

  location = /akamai/sureroute-test-object.html {
    autoindex off;
    add_header Cache-Control public;
    expires 4w;
    try_files $uri $uri/ =404;
  }

  location = /favicon.ico {
    return 301 https://your_domain.com/favicon.ico;
  }

  location / {
    proxy_pass   $boilderplatenode;
  }
}