package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	_db *sqlx.DB
}

var postgreSQL *PostgreSQL

func initPostgreSQL() {
	// TODO: init postgreSQL
	fmt.Println("implement me init postgreSQL")
}

func (p PostgreSQL) GetDB() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQL) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQL) Close() error {
	//TODO implement me
	panic("implement me")
}
