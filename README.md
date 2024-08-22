# Introduction 
Use go-boilerplate to start a http server project in Golang

# What is included
1. Configuration driven server
2. DB connect with multiple driver support
3. Daemon script
4. logrotate config
5. nginx proxy
6. html template
7. panic and recovery

# Repository Structure
1. app.go at root is main file
2. files/ contains all production setup configuration
3. files/var/www contains assets, html and templates

```
   |-api
   |-common
   |---config
   |---constant
   |---database
   |---slack
   |---utils
   |-files
   |---etc
   |-----go-boilerplate
   |-------development
   |-------production
   |-------staging
   |-----init
   |-----logrotate.d
   |-----nginx
   |-------sites-available
   |---setup
   |---sql
   |---var
   |-----www
   |-------email_templates
   |-------go-boilerplate
   |---------html
   |---------public
   |-----------media
   |-------------css
   |-------------img
   |-------------script
   |-----------template
   |-------------base
   |-------------header
   |---------------mobile
   |---------------web
   |-lib
   |-model
   |-tools
   |---panics
   |-----example
   |---safe
   |-type
   |-app.go

```

#Setup
1. go mod tidy
2. go mod vendor
3. go run app.go
4. Update DB config to connect with right DB. 
5. Replate YOUR_ENV by production env definition key
6. Rename files/etc/init/org-go-boilerpate-cron.conf to files/etc/init/<your-org-name>-<repo-name>-cron.conf
7. Rename files/etc/init/org-go-boilerpate.conf to files/etc/init/<your-org-name>-<repo-name>.conf




# How to use this Repository
1. git clone git@github.com:Somesh/go-boilerplate
2. chmod x setup.sh
3. ./setup.sh "go-boilerplate" "new-repo" "org_name"
