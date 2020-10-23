package config

import (
	"Nani/internal/app/file"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

//Application config
type Config struct {
	ApiUrl    string `yaml:api_url`
	Hl        string `yaml:hl`
	Gl        string `yaml:gl`
	Key       string
	KeysCount int
	AppsCount int

	envs []string `yaml: ,flow`
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
func New(p ...string) *Config {
	path := "/config/dev.yml"
	if len(p) > 0 {
		path = p[0]
	}

	data, err := file.ReadAll(path)
	if err != nil {
		log.Fatal(err)
	}

	var config *Config
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		log.Fatal(err)
	}

	envs := loadEnvs(config.envs...)

	v, ok := envs["api_key"]
	if !ok {
		log.Print("Key api_key not found in sys envs")
	} else {
		config.Key = v
	}

	return config
}
