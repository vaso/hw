package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

type EnvConfig struct {
	Grpc struct {
		Port int `envconfig:"RATE_GRPC_PORT" default:"80"`
	}
	DB struct {
		Host string `envconfig:"RATE_DB_HOST"`
		User string `envconfig:"RATE_DB_USER"`
		Pass string `envconfig:"RATE_DB_PASSWORD"`
		Name string `envconfig:"RATE_DB_NAME"`
		Port int    `envconfig:"RATE_DB_PORT"`
	}
	RedisPort int `envconfig:"RATE_REDIS_PORT"`
	K         int `envconfig:"K"`
	M         int `envconfig:"M"`
	N         int `envconfig:"N"`
}

func GetEnvConfig(configFile string) *EnvConfig {
	var cfg EnvConfig
	err := gotenv.Load(configFile)
	if err != nil {
		panic(err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
