package model

type HealthCheckResult []int64

func HealthCheck(slave bool) error {
	var result HealthCheckResult
	var err error

	if !slave {
		err = statements.masterHealth.Select(&result)
	} else {
		err = statements.slaveHealth.Select(&result)

	}

	return err
}
