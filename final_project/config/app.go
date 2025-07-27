package config

type AppConfig struct {
	K           int
	M           int
	N           int
	IPWhitelist []string
	IPBlacklist []string
	EnvConfig   EnvConfig
}

func GetAppConfig(envConfig EnvConfig) *AppConfig {
	return &AppConfig{
		K:         envConfig.K,
		M:         envConfig.M,
		N:         envConfig.N,
		EnvConfig: envConfig,
	}
}
