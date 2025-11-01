package config

import "time"

// Configuration management

type Config struct {
	Env      *Env      // Environment variables configuration
	Database *Database // Database configuration
}

func LoadConfig() (*Config, error) {
	// Load environment variables
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	// Load database configuration
	database, err := NewMongoDatabase(MongoDBConfig{
		URI:        env.MongoURI,
		Database:   env.MongoDatabase,
		MaxRetries: 4,
		RetryDelay: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Config{
		Env: env,
		Database: database,
	}, nil
}
