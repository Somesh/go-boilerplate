package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/gops/agent"

	"github.com/tokopedia/grace"
	"gopkg.in/tokopedia/logging.v1"
	"gopkg.in/tokopedia/logging.v1/tracer"

	//inits
	"github.com/Somesh/go-boilerplate/api"
	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/common/database"

	"github.com/Somesh/go-boilerplate/src/manager"
	// Azure Eventhub
	//	_ "github.com/Somesh/go-boilerplate/event/eventhub/publisher"
	// NSQ
	_ "github.com/Somesh/go-boilerplate/event/nsq/publisher"
)

// TODO Add prefixes to topic names before publishing
func main() {
	if err := startApp(false); err != nil {
		log.Printf("Unable To Start App Error : %+v", err)
	}
}

func startApp(isTest bool) error {

	var configtest = flag.Bool("test", false, "config test")

	flag.Parse()

	if !isTest {
		logging.LogInit()
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logging.Debug.Println("app started") // message will not appear unless run with -debug switch

	//gops helps us get stack trace if something wrong/slow in production
	opts := agent.Options{
		ShutdownCleanup: true,
	}

	if err := agent.Listen(opts); err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfig()

	if isTest {
		//TODO :: For Unit test there is no need for below steps.
		return nil
	}

	database.Init(cfg)

	// //Initialise Web Service and HTTP Handler
	api := api.InitAPIMod(cfg)
	api.InitHandlers()

	managerMod := manager.New()
	managerMod.Init()

	if !logging.IsDebug() {
		go logging.StatsLogInterval(5, true)
	}

	if *configtest {
		os.Exit(0)
	}

	tracer.Init(&cfg.Tracer)
	log.Fatal(grace.ServeWithConfig(":9000", cfg.Grace.ToGraceConfig(), nil))

	return nil
}
