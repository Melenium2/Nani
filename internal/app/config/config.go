package config

import (
	"Nani/internal/app/file"
	"database/sql"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type DBConfig struct {
	Database   string `yaml:"name"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Address    string `yaml:"address"`
	Port       string `yaml:"port"`
	Schema     string `yaml:"schema"`
	Connection *sql.DB
}

//Application config
type Config struct {
	ApiUrl    string   `yaml:"api_url"`
	Hl        string   `yaml:"hl"`
	Gl        string   `yaml:"gl"`
	Database  DBConfig `yaml:"database"`
	Key       string
	KeysCount int
	AppsCount int

	Envs []string `yaml:",flow"`
}

/**
Load system environment variables from given array
*/
func loadEnvs(e ...string) map[string]string {
	envs := make(map[string]string)

	for _, k := range e {
		envs[k] = os.Getenv(k)
	}

	return envs
}

/**
Create new instance of app config
@Param p : String (path to custom config file.yml)
*/
func New(p ...string) Config {
	path := "../../../config/dev.yml"
	if len(p) > 0 {
		path = p[0]
	}

	data, err := file.ReadAll(path)
	if err != nil {
		panic(err)
	}

	config := Config{}
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		panic(err)
	}

	envs := loadEnvs(config.Envs...)

	v, ok := envs["api_key"]
	if !ok {
		log.Print("Key api_key not found in sys envs")
	} else {
		config.Key = v
	}
	v, ok = envs["db_pass"]
	if ok {
		config.Database.Password = v
	}
	v, ok = envs["db_user"]
	if ok {
		config.Database.User = v
	}

	return config
}
