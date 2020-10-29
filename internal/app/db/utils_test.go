package db_test

import (
	config2 "Nani/internal/app/config"
	"Nani/internal/app/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionUrl_ShouldReturnRightUrl_NoErrors(t *testing.T) {
	config := config2.DBConfig{
		Database: "default",
		User:     "admin",
		Password: "123456",
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}

	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	assert.Equal(t, "http://admin:123456@192.168.99.100:8123/default", url)
}

func TestConnectionUrl_ShouldReturnUrlWithoutUserAndPass_NoError(t *testing.T) {
	config := config2.DBConfig{
		Database: "default",
		User:     "admin",
		Password: "",
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}

	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	assert.Equal(t, "http://192.168.99.100:8123/default", url)
}

func TestConnectionUrl_ShouldReturnErrorCozAddressIsEmpty_Error(t *testing.T) {
	config := config2.DBConfig{
		Database: "default",
		User:     "admin",
		Password: "123456",
		Address:  "",
		Port:     "8123",
		Schema:   "",
	}

	url, err := db.ConnectionUrl(config)
	assert.Error(t, err)
	assert.Equal(t, "empty db address", err.Error())
	assert.Empty(t, url)
}

func TestConnectionUrl_ShouldReturnErrorCozPortIsEmpty_Error(t *testing.T) {
	config := config2.DBConfig{
		Database: "default",
		User:     "admin",
		Password: "123456",
		Address:  "192.168.99.100",
		Port:     "",
		Schema:   "",
	}

	url, err := db.ConnectionUrl(config)
	assert.Error(t, err)
	assert.Equal(t, "empty db port", err.Error())
	assert.Empty(t, url)
}

func TestConnectionUrl_ShouldReturnDefaultDbIfNoPresented_NoError(t *testing.T) {
	config := config2.DBConfig{
		Database: "",
		User:     "admin",
		Password: "123456",
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}

	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	assert.Equal(t, "http://admin:123456@192.168.99.100:8123/default", url)
}

func TestConnect_ShouldEstablishNewConnectionToDatabase_NoError(t *testing.T) {
	config := config2.DBConfig{
		Password: "123456",
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}
	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	conn, err := db.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestConnect_ShouldReturnErrorCozDbCanNotEstablishConnection_NoError(t *testing.T) {
	config := config2.DBConfig{
		Address:  "192.168.99.101",
		Port:     "8123",
		Schema:   "",
	}
	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	conn, err := db.Connect(url)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestInitSchema_ShouldCreateNewTableInDatabase_NoError(t *testing.T) {
	config := config2.DBConfig{
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}
	url, err := db.ConnectionUrl(config)
	assert.NoError(t, err)
	conn, err := db.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	err = db.InitSchema(conn, "../../../config/schema.sql")
	assert.NoError(t, err)

	_, err = conn.Exec(`select * from apps`)
	assert.NoError(t, err)
	_, err = conn.Exec("drop table apps")
	assert.NoError(t, err)
}

