package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Somesh/go-boilerplate/common/config"

	"github.com/Somesh/go-boilerplate/model"
	"gopkg.in/tokopedia/logging.v1"
)

// individual error thresholds for each type of health check
var errorCounter map[string]int = make(map[string]int)

// Health perform health checks
func (m *APIMod) Health(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	// ctx := context.Background()

	// SS:  Use when detail health check is needed
	// if errHealthCheck := raiseAlarm(ctx, m.doHealthChecks(ctx)); errHealthCheck != nil {
	// 	log.Printf("[api][Health] Error : %+v", errHealthCheck)
	// 	go postToSlack(ctx, errHealthCheck, m.log)
	// 	return nil, errHealthCheck
	// }

	return "OK", nil
}

func (m *APIMod) doHealthChecks(ctx context.Context) error {
	defer timeTrack(time.Now(), "doHealthChecks")
	var errHealthCheck error
	// Checking DB Master and Slave. We will raise an alarm only if
	// number of successive db errors cross a threshold
	if errHealthCheck = m.dbChecks(ctx); errHealthCheck != nil {
		return errHealthCheck
	}

	// Checking Redis. We will raise an alarm only if
	// number of successive redis errors cross a threshold
	if errHealthCheck = m.redisHealthCheck(ctx); errHealthCheck != nil {
		return errHealthCheck
	}

	return nil
}

func (m *APIMod) redisHealthCheck(ctx context.Context) error {
	return nil
}

func (m *APIMod) dbChecks(ctx context.Context) error {
	defer timeTrack(time.Now(), "dbChecks")
	var errDb error
	// Checking DB Master
	if errDb = model.HealthCheck(false); errDb != nil {
		return fmt.Errorf("DB Master Down. Error : {%+v}", errDb)
	}

	// Checking DB Slave
	if errDb = model.HealthCheck(true); errDb != nil {
		return fmt.Errorf("DB Slave Down. Error : {%+v}", errDb)
	}

	return nil
}

func postToSlack(ctx context.Context, err error, slackWriter *log.Logger) {
	defer timeTrack(time.Now(), "postToSlack")
	if slackWriter == nil || err == nil {
		return
	}

	ip := config.GetLocalIP()
	env := os.Getenv("YOUR_ENV")
	slackWriter.Println(fmt.Sprintf("[%s][%s] %s", ip, env, err.Error()))
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	logging.Debug.Printf("%s took %s", name, elapsed)
}

// Increments a counter in case of error and raises error only if counter > threshold
func raiseAlarm(ctx context.Context, err error) error {
	if err == nil {
		errorCounter = map[string]int{}
		return nil
	}
	errStr := err.Error()
	threshold := config.GetConfig().Server.FailHealthCheckThreshold
	if errorCounter[errStr] += 1; errorCounter[errStr] > threshold {
		return err
	}
	return nil
}
