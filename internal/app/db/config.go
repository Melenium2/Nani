package db

import "database/sql"

type Config struct {
	Database   string
	User       string
	Password   string
	Address    string
	Port       string
	Schema     string
	Connection *sql.DB
}
