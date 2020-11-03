package db

import (
	config2 "Nani/internal/app/config"
	"Nani/internal/app/file"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func ConnectionUrl(config config2.DBConfig) (string, error) {
	url := "http://"

	if config.User != "" && config.Password != "" {
		url += config.User + ":" + config.Password + "@"
	}

	if config.Address == "" {
		return "", errors.New("empty db address")
	}
	url += config.Address

	if config.Port == "" {
		return "", errors.New("empty db port")
	}
	url += ":" + config.Port + "/"

	db := config.Database
	if db == "" {
		db = "default"
	}
	url += db

	return url, nil
}

func Connect(url string) (*sql.DB, error) {
	connect, err := sql.Open("clickhouse", url)
	if err != nil {
		return nil, err
	}

	if err := connect.Ping(); err != nil {
		return nil, err
	}

	connect.SetConnMaxLifetime(time.Second * 30)
	connect.SetMaxOpenConns(10)

	return connect, nil
}

func InitSchema(connection *sql.DB, schemafile string) error {
	f := file.New(schemafile)
	b, err := f.ReadAll()
	if err != nil {
		return err
	}
	schema := string(b)
	tables := strings.Split(schema, ";\n")

	for _, v := range tables {
		_, err = connection.Exec(v)
		if err != nil {
			if strings.Contains(err.Error(), "Code: 57") {
				newsql := strings.ReplaceAll(v,"CREATE TABLE", "ATTACH TABLE")
				if _, err = connection.Exec(newsql); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}


	return nil
}
