package model

import (
	"github.com/jmoiron/sqlx"

	"github.com/Somesh/go-boilerplate/common/database"
)

// this is done for the time being
var db *database.MasterSlave

type PreparedStatements struct {
	masterHealth *sqlx.Stmt
	slaveHealth  *sqlx.Stmt
}

var (
	statements PreparedStatements
)

func Init(dbPtr *database.MasterSlave) {

	db = dbPtr

	statements.masterHealth = db.Master.Preparex(healthQuery)
	statements.slaveHealth = db.Slave.Preparex(healthQuery)

	// init item action rules
	// Load rows required to init service
}

var (
	healthQuery = `/* ping */ SELECT 1;`
)
