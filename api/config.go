package api

import (

	"net/http"

	"gopkg.in/tokopedia/logging.v1"
)


// Health perform health checks
func (m *APIMod) ConfigHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {


	logging.Debug.Printf("%+v", m.cfg)
	// ctx := context.Background()

	// SS:  Use when detail health check is needed
	// if errHealthCheck := raiseAlarm(ctx, m.doHealthChecks(ctx)); errHealthCheck != nil {
	// 	log.Printf("[api][Health] Error : %+v", errHealthCheck)
	// 	go postToSlack(ctx, errHealthCheck, m.log)
	// 	return nil, errHealthCheck
	// }

	return m.cfg, nil
}