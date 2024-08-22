package slack

// import (
// 	"io/ioutil"
// 	"log"

// 	"github.com/Somesh/go-boilerplate/common/config"
// 	"github.com/Somesh/go-boilerplate/common/constant"

// 	slacker "github.com/whosonfirst/go-writer-slackcat"
// )

// var (
// 	SlackLogger         *log.Logger
// 	SlackDebugLogger    *log.Logger
// 	SlackBusinessLogger map[string]*log.Logger
// )

// func Init(cfgs *config.Config) {
// 	if cfgs.Slack.WebhookUrl != "" {
// 		w := slacker.Writer{
// 			Config: &cfgs.Slack,
// 		}
// 		SlackLogger = log.New(w, "", log.Ldate|log.Ltime)
// 	} else {
// 		SlackLogger = log.New(ioutil.Discard, "", 0)
// 	}

// 	if cfgs.DebugSlack.WebhookUrl != "" {
// 		w := slacker.Writer{
// 			Config: &cfgs.DebugSlack,
// 		}
// 		SlackDebugLogger = log.New(w, "", log.Ldate|log.Ltime)
// 	} else {
// 		SlackDebugLogger = log.New(ioutil.Discard, "", 0)
// 	}

// 	// initiating slack for specific businesses
// 	SlackBusinessLogger = make(map[string]*log.Logger)
// 	if cfgs.BusinessSlack.WebhookUrl != "" {
// 		w := slacker.Writer{
// 			Config: &cfgs.BusinessSlack,
// 		}

// 		SlackBusinessLogger[constant.CommonString] = log.New(w, "", log.Ldate|log.Ltime)

// 		if cfgs.BusinessSlackCustomChannel.EventsChannel != "" {
// 			eventsBusinessSlackConfig := cfgs.BusinessSlack
// 			eventsBusinessSlackConfig.Channel = cfgs.BusinessSlackCustomChannel.EventsChannel
// 			writerEvents := slacker.Writer{
// 				Config: &eventsBusinessSlackConfig,
// 			}
// 			SlackBusinessLogger[constant.EventsName] = log.New(writerEvents, "", log.Ldate|log.Ltime)
// 		} else {
// 			SlackBusinessLogger[constant.EventsName] = log.New(w, "", log.Ldate|log.Ltime)
// 		}

// 		if cfgs.BusinessSlackCustomChannel.DealsChannel != "" {
// 			dealsBusinessSlackConfig := cfgs.BusinessSlack
// 			dealsBusinessSlackConfig.Channel = cfgs.BusinessSlackCustomChannel.DealsChannel
// 			writerDeals := slacker.Writer{
// 				Config: &dealsBusinessSlackConfig,
// 			}
// 			SlackBusinessLogger[constant.DealsName] = log.New(writerDeals, "", log.Ldate|log.Ltime)
// 		} else {
// 			SlackBusinessLogger[constant.DealsName] = log.New(w, "", log.Ldate|log.Ltime)
// 		}
// 	} else {
// 		SlackBusinessLogger[constant.CommonString] = log.New(ioutil.Discard, "", 0)
// 		SlackBusinessLogger[constant.EventsName] = log.New(ioutil.Discard, "", 0)
// 		SlackBusinessLogger[constant.DealsName] = log.New(ioutil.Discard, "", 0)
// 	}
// }

// func GetLogger() *log.Logger {
// 	return SlackLogger
// }

// func GetDebugLogger() *log.Logger {
// 	return SlackDebugLogger
// }

// func GetBusinessLogger() map[string]*log.Logger {
// 	return SlackBusinessLogger
// }
