package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tokopedia/grace"
	slacker "github.com/whosonfirst/go-writer-slackcat"
	"github.com/whosonfirst/slackcat"
	"gopkg.in/gcfg.v1"
	"gopkg.in/tokopedia/logging.v1/tracer"

	"github.com/Somesh/go-boilerplate/common/constant"
	"github.com/Somesh/go-boilerplate/lib"
)

var CF *Config

type Config struct {
	Server             ServerConfig
	Database           map[string]*DatabaseConfig
	DatabaseConnection DatabaseConnectionConfig
	Tracer             tracer.Config
	Timeout            time.Duration
	Slack              slackcat.Config
	Grace              GraceCfg
	Event              EventCfg
	ConnectionHub      map[string]*ConnectionHubCfg
	NSQ                NSQCfg
}

type DatabaseConfig struct {
	Master        string `json:"master"`
	Slave         string `json:"slave"`
	MaxMasterConn int
	MaxSlaveConn  int
	Driver        string
}

type DatabaseConnectionConfig struct {
	MasterMaxOpenConn int
	MasterMaxIdleConn int
	SlaveMaxOpenConn  int
	SlaveMaxIdleConn  int
	MaxOpenLifetime   int
	MaxIdleLifetime   int
	PingRetryInterval int
}

type ServerConfig struct {
	Host                     string
	Port                     int
	TemplatePath             string
	Timeout                  time.Duration
	PingClientTimeout        time.Duration
	LocalIP                  string
	FailHealthCheckThreshold int
}

type GraceCfg struct {
	Timeout          string
	HTTPReadTimeout  string
	HTTPWriteTimeout string
}

type EventCfg struct {
	APIURL        string
	HubName       string
	HubConnString string
	HubNameSpace  string
}

type ConnectionHubCfg struct {
	HubName       string
	HubConnString string
	Partitions    []string
}

type NSQCfg struct {
	ListenAddress  []string
	LookUpAddress  []string
	PublishAddress string
	Prefix         string
}

func (g GraceCfg) ToGraceConfig() grace.Config {
	timeout, err := time.ParseDuration(g.Timeout)
	if err != nil {
		timeout = time.Second * 5
	}

	readTimeout, err := time.ParseDuration(g.HTTPReadTimeout)
	if err != nil {
		readTimeout = time.Second * 10
	}

	writeTimeout, err := time.ParseDuration(g.HTTPWriteTimeout)
	if err != nil {
		writeTimeout = time.Second * 10
	}

	return grace.Config{
		Timeout:          timeout,
		HTTPReadTimeout:  readTimeout,
		HTTPWriteTimeout: writeTimeout,
	}
}

func init() {
	CF = &Config{}
	GOPATH := os.Getenv("GOPATH")
	ok := ReadConfig(CF, "/etc", "go-boilerplate") || ReadConfig(CF, GOPATH+"/src/github.com/Somesh/go-boilerplate/files/etc", "go-boilerplate") || ReadConfig(CF, "files/etc", "go-boilerplate")
	if !ok {
		log.Fatal("Failed to read config file")
	}
	SetLocalIP()
}

// ReadConfig is file handler for reading configuration files into variable
// Param: - config pointer of Config, filepath string
// Return: - boolean
func ReadConfig(cfg *Config, path string, module string) bool {
	environ := os.Getenv("YOUR_ENV")
	if environ == "" {
		environ = constant.ENV_DEVELOPMENT
	}

	environ = strings.ToLower(environ)

	parts := []string{"main"}
	var configString []string

	for _, v := range parts {
		fname := path + "/" + module + "/" + environ + "/" + module + "." + v + ".ini"
		fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "Reading", fname)

		config, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Println("common/config.go function ReadConfig", err)
			return false
		}

		configString = append(configString, string(config))
	}

	err := gcfg.ReadStringInto(cfg, strings.Join(configString, "\n\n"))
	if err != nil {
		log.Println("common/config.go function ReadConfig", err)
		return false
	}

	return true
}

func GetConfig() *Config {
	return CF
}

func GetLogger() *log.Logger {
	var slog *log.Logger

	if CF.Slack.WebhookUrl != "" {
		w := slacker.Writer{
			Config: &CF.Slack,
		}
		slog = log.New(w, "", log.Ldate|log.Ltime)
	} else {
		slog = log.New(ioutil.Discard, "", 0)
	}
	return slog
}

func GetLocalIP() string {
	return CF.Server.LocalIP
}

func SetLocalIP() {
	// Once Set , we do not want to Reset it again
	if CF.Server.LocalIP != "" {
		return
	}
	localIP, err := lib.GetOutboundIP()
	if err != nil {
		log.Println("[common/config.go][SetLocalIP] Error Fetching IP. Error : ", err)
	} else {
		CF.Server.LocalIP = localIP
	}
}
