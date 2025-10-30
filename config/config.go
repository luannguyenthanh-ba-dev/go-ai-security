package config

// Configuration management

type Config struct {
	Env *Env // Environment variables configuration
}

func LoadConfig() (*Config, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		Env: env,
	}, nil
}
