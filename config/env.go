package config

// Environment configuration
import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Environment variables configuration
// Struct is as class in other languages like Python, Java, etc.
// Tags are used to map the environment variables to the struct fields
type Env struct {
	AppName                string `mapstructure:"APP_NAME"`
	Port                   string `mapstructure:"PORT"`
	AppEnv                 string `mapstructure:"APP_ENV"`
	MongoURI               string `mapstructure:"MONGO_URI"`
	MongoDatabase          string `mapstructure:"MONGO_DATABASE"`
	RedisURL               string `mapstructure:"REDIS_URL"`
	EmailResendAPIKey      string `mapstructure:"EMAIL_RESEND_API_KEY"`
	PasswordHashSaltRounds int    `mapstructure:"PASSWORD_HASH_SALT_ROUNDS"`
	JWTSecret              string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn           int    `mapstructure:"JWT_EXPIRES_IN"`
}

// Return *Env and error: *Env is the environment variables configuration, error is the error if any
func LoadEnv() (*Env, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		// Log the error using the logger zap - L function is used to get the logger instance
		zap.L().Error("Failed to read config file", zap.Error(err))
		return nil, err
	}

	// Unmarshal the environment variables into the Env struct
	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		zap.L().Error("Failed to unmarshal config file", zap.Error(err))
		return nil, err
	}

	return &env, nil
}
