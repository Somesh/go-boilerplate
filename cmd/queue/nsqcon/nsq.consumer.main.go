package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/google/gops/agent"
	"gopkg.in/tokopedia/logging.v1"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/common/database"
	"github.com/Somesh/go-boilerplate/src/manager"

	"github.com/Somesh/go-boilerplate/event/nsq"
	nsqConsumer "github.com/Somesh/go-boilerplate/event/nsq/consumer"
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

	wg := &sync.WaitGroup{}
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
	// api := api.InitAPIMod(cfg)
	// api.InitHandlers()

	if !logging.IsDebug() {
		go logging.StatsLogInterval(5, true)
	}

	if !*configtest {
		lookupt := flag.Lookup("t")
		if lookupt != nil {
			tFlag, _ := strconv.ParseBool(lookupt.Value.String())
			configtest = &tFlag
		}
	}

	if *configtest {
		os.Exit(0)
	}

	nsqcfg := nsqrun.Options{
		ListenAddress:  cfg.NSQ.ListenAddress,
		PublishAddress: cfg.NSQ.PublishAddress,
		Prefix:         cfg.NSQ.Prefix,
		LookUpAddress:  cfg.NSQ.LookUpAddress,
	}

	managerMod := manager.New()
	managerMod.Init()

	cons := nsqConsumer.Setup(managerMod, cfg)
	queueServer := nsqrun.New(&nsqcfg, cons)

	if isTest {
		return nil
	}

	wg.Add(1)
	go queueServer.Run()

	if !logging.IsDebug() {
		wg.Add(1)
		go logging.StatsLogInterval(5, true)
	}

	if *configtest {
		os.Exit(0)
	}

	wg.Wait()
	return nil
}
