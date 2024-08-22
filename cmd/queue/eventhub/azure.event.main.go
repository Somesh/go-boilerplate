package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/google/gops/agent"
	"gopkg.in/tokopedia/logging.v1"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/common/database"
	"github.com/Somesh/go-boilerplate/src/manager"

	"github.com/Somesh/go-boilerplate/event/eventhub"
	"github.com/Somesh/go-boilerplate/event/eventhub/consumer"
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

	eventOpts := event.Options{
		HubConnString: cfg.Event.HubConnString,
		HubName:       cfg.Event.HubName,
		HubNameSpace:  cfg.Event.HubNameSpace,
	}

	managerMod := manager.New()
	managerMod.Init()

	cons := consumer.Setup(managerMod, cfg)
	queueServer := event.New(&eventOpts, cons, cfg)

	queueServer.Run()

	return nil
}

// setupSignalHandler sets up a signal channel for handling OS interrupts, waits for the signal, and then gracefully shuts down.
func setupSignalHandler(cancel context.CancelFunc, wg *sync.WaitGroup) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal for graceful shutdown
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, shutting down gracefully...")

		// Initiate graceful shutdown
		cancel()
		wg.Wait()
		fmt.Println("Shutdown complete.")
	}()
}
