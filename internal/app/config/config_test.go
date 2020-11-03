package config_test

import (
	"Nani/internal/app/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew_ShouldCreateConfigInstanceWithAllFieldsFromConfigFile_NoError(t *testing.T) {
	c := config.New()
	assert.NotEmpty(t, c.ApiUrl)
	assert.NotEmpty(t, c.Gl)
	assert.NotEmpty(t, c.Hl)
	assert.NotEmpty(t, c.Database.Schema)
	assert.NotEmpty(t, c.Database.Address)
	assert.NotEmpty(t, c.Database.Database)
	assert.NotEmpty(t, c.Database.Port)
}

func TestNew_ShouldCreateConfigInstanceWithEnvVars_NoError(t *testing.T) {
	os.Setenv("api_key", "123")
	c := config.New()
	assert.NotEmpty(t, c.Key)
	assert.Equal(t, "123", c.Key)
}
