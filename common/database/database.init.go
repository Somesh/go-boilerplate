package database

import (

	// "context"
	"database/sql"
	"log"
	"time"

	"gopkg.in/tokopedia/logging.v1"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	config "github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/tools/safe"
	// slack 	"github.com/Somesh/go-boilerplate/common/slack"
)

type MasterSlave struct {
	Master *DB
	Slave  *DB
}

var DBConnMap map[string]*MasterSlave

//TODO: Set prometheus

// DB configuration
type DB struct {
	DBConnection    *sqlx.DB
	DBString        string
	RetryInterval   int
	MaxOpenConn     int
	MaxIdleConn     int
	MaxOpenLifetime int
	MaxIdleLifetime int
	doneChannel     chan bool
}

var (
	dbTicker *time.Ticker
)

func Init(cfgs *config.Config) {
	DBConnMap = make(map[string]*MasterSlave)

	for k, cfg := range cfgs.Database {
		masterDsn := cfg.Master
		slaveDsn := cfg.Slave
		driver := cfg.Driver
		var maxConnSlave, maxConnMaster, retryInterval int
		if maxConnMaster = cfgs.DatabaseConnection.MasterMaxOpenConn; maxConnMaster == 0 {
			maxConnMaster = 80
		}
		if maxConnSlave = cfgs.DatabaseConnection.SlaveMaxOpenConn; maxConnSlave == 0 {
			maxConnSlave = 100
		}
		if retryInterval = cfgs.DatabaseConnection.PingRetryInterval; retryInterval == 0 {
			retryInterval = 10
		}

		Master := &DB{
			DBString:        masterDsn,
			RetryInterval:   retryInterval,
			MaxOpenConn:     maxConnMaster,
			MaxIdleConn:     cfgs.DatabaseConnection.MasterMaxIdleConn,
			MaxOpenLifetime: cfgs.DatabaseConnection.MaxOpenLifetime,
			MaxIdleLifetime: cfgs.DatabaseConnection.MaxIdleLifetime,
			doneChannel:     make(chan bool),
		}
		Master.ConnectAndMonitor(driver, k+" master")

		Slave := &DB{
			DBString:        slaveDsn,
			RetryInterval:   retryInterval,
			MaxOpenConn:     maxConnSlave,
			MaxIdleConn:     cfgs.DatabaseConnection.SlaveMaxIdleConn,
			MaxOpenLifetime: cfgs.DatabaseConnection.MaxOpenLifetime,
			MaxIdleLifetime: cfgs.DatabaseConnection.MaxIdleLifetime,
			doneChannel:     make(chan bool),
		}
		Slave.ConnectAndMonitor(driver, k+" slave")

		DBConnMap[k] = &MasterSlave{
			Master: Master,
			Slave:  Slave,
		}

	}

	dbTicker = time.NewTicker(time.Second * 2)
}

func (d *DB) GetDB() *sqlx.DB {
	return d.DBConnection
}

// Connect to database
func (d *DB) Connect(driver string) error {
	var db *sqlx.DB
	var err error

	db, err = sqlx.Open(driver, d.DBString)
	if err != nil {
		// tdklog.StdError(context.Background(), nil, err, "[Error]: DB open connection error")
		logging.Debug.Println(err, "[Error]: DB open connection error")
		return err
	}

	if d.MaxOpenConn > 0 {
		db.SetMaxOpenConns(d.MaxOpenConn)
	}
	if d.MaxIdleConn > 0 {
		db.SetMaxIdleConns(d.MaxIdleConn)
	}
	if d.MaxOpenLifetime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(d.MaxOpenLifetime))
	}
	// if d.MaxIdleLifetime > 0 {
	// 	db.SetConnMaxIdleTime(time.Second * time.Duration(d.MaxIdleLifetime))
	// }

	d.DBConnection = db

	return nil
}

// ConnectAndMonitor to database
func (d *DB) ConnectAndMonitor(driver, name string) {
	err := d.Connect(driver)

	// ctx := context.Background()
	log.Println(driver, name, "Error: ", err)
	if err != nil {
		// tdklog.StdErrorf(ctx, nil, err, "Not connected to database %s, trying", name)
		log.Fatal("Error: %+v, Not connected to database %s, trying", err, name)
	} else {
		// tdklog.StdDebugf(ctx, nil, nil, "Success connecting to database %s", name)
		logging.Debug.Printf("Success connecting to database %s", name)
	}

	ticker := time.NewTicker(time.Duration(d.RetryInterval) * time.Second)
	go func() {
		defer safe.Recover()
		for {
			select {
			case <-ticker.C:
				if d.DBConnection == nil {
					d.Connect(driver)
				} else {
					err := d.DBConnection.Ping()
					if err != nil {
						// tdklog.StdErrorf(ctx, nil, err, "[Error]: DB reconnect to database %s error", name)
						logging.Debug.Printf("Error: %+v, Not connected to database %s, trying", err, name)
					}
				}
			case <-d.doneChannel:
				return
			}
		}
	}()
}

// DoneConnectAndMonitor to exit connect and monitor
func (d *DB) DoneConnectAndMonitor() {
	d.doneChannel <- true
}

// Prepare query for database queries
func (d *DB) Prepare(query string) *sql.Stmt {
	statement, err := d.DBConnection.Prepare(query)

	if err != nil {
		// tdklog.StdError(context.Background(), query, err, "Failed to prepare query")
		log.Fatalf("Query: %+v, Failed to prepare query. Error:  %s,", query, err)
		// slack.GetLogger().Printf("Failed to prepare query: %s. Error: %s", query, err.Error())
	}

	return statement
}

// Preparex query for database queries
func (d *DB) Preparex(query string) *sqlx.Stmt {
	if d == nil {
		// tdklog.StdError(context.Background(), query, nil, "Failed to preparex query")
		log.Fatalf("Query: %+v, Failed to prepare query. Db Pointer is nil", query)
		// slack.GetLogger().Printf("Failed to preparex query: %s. Error: db is nil", query)
		return nil
	}

	statement, err := d.DBConnection.Preparex(query)

	if err != nil {
		// tdklog.StdError(context.Background(), query, err, "Failed to preparex query")
		log.Fatalf("Query: %+v, Failed to prepare query. Error:  %s,", query, err)
		// slack.GetLogger().Printf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

// Transactionx is Wrapper for db transaction, make database operations (SELECT, INSERT, UPDATE, etc.) into a single unit.
// Note: in task function, the content should for DB tasks only to avoid long-running query and occupying the connection
// pool for extended periods.
func (d *DB) Transactionx(taskName string, task func(tc *sqlx.Tx) error) error {
	if d == nil {
		// tdklog.StdError(context.Background(), taskName, nil, "Failed to do transactionx")
		log.Fatalf("taskName: %+v, Failed to do transactionx. Db Pointer is nil", taskName)
		// slack.GetLogger().Printf("Failed to  to do transactionx: %s. Error: db is nil", taskName)
		return nil
	}

	tx, err := d.DBConnection.Beginx()
	if err != nil {
		return err
	}

	err = task(tx)
	if err != nil {
		// tdklog.StdError(context.Background(), taskName, err, "Failed to run transaction, rollback is executed")
		logging.Debug.Printf("taskName: %+v, Failed to run transaction, rollback is executed. Error:  %s,", taskName, err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			// tdklog.StdError(context.Background(), taskName, errRollback, "Failed to rollback")
			logging.Debug.Printf("taskName: %+v, Failed to rollback. Error:  %s,", taskName, err)
		}

		return err
	}

	return tx.Commit()
}
