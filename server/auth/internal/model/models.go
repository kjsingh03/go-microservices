package model

import (
	"database/sql"
)

var db *sql.DB

func SetDB(dbPool *sql.DB) {
	db = dbPool
}