package config

// Database configuration

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBConfig struct {
	URI        string
	Database   string
	MaxRetries int
	RetryDelay time.Duration
}

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDatabase creates a new MongoDB database connection with retry mechanism
func NewMongoDatabase(mongoConfig MongoDBConfig) (*Database, error) {
	// Set default values if not provided
	maxRetries := mongoConfig.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // Default 3 retries
	}
	retryDelay := mongoConfig.RetryDelay
	if retryDelay == 0 {
		retryDelay = 2 * time.Second // Default 2 seconds
	}

	// Create MongoDB client object (this doesn't actually connect to server)
	// Connection happens lazily on first operation (e.g., Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	cancel()

	if err != nil {
		zap.L().Error("failed to create MongoDB client",
			zap.String("uri", mongoConfig.URI),
			zap.Error(err),
		)
		return nil, err
	}

	zap.L().Debug("MongoDB client object created",
		zap.String("uri", mongoConfig.URI),
	)

	// Retry Ping to verify actual connection (this is where real connection happens)
	currentDelay := retryDelay
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			zap.L().Info("retrying MongoDB connection (ping)",
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
				zap.Duration("delay", currentDelay),
			)
			time.Sleep(currentDelay)
			// Exponential backoff: double the delay for each retry
			currentDelay *= 2
		}

		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Ping(pingCtx, nil)
		pingCancel()

		if err == nil {
			zap.L().Info("successfully connected to MongoDB",
				zap.String("uri", mongoConfig.URI),
				zap.String("database", mongoConfig.Database),
				zap.Int("total_attempts", attempt+1),
			)
			break
		}

		zap.L().Warn("failed to connect to MongoDB (ping failed)",
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.String("database", mongoConfig.Database),
			zap.Error(err),
		)

		// If this is the last attempt, cleanup and return error
		if attempt == maxRetries {
			zap.L().Error("failed to connect to MongoDB after all retries",
				zap.Int("total_attempts", maxRetries+1),
				zap.String("uri", mongoConfig.URI),
				zap.String("database", mongoConfig.Database),
				zap.Error(err),
			)
			// Cleanup: disconnect client before returning error
			disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer disconnectCancel()
			_ = client.Disconnect(disconnectCtx)
			return nil, err
		}
	}

	return &Database{
		Client:   client,
		Database: client.Database(mongoConfig.Database),
	}, nil
}

// Close closes the MongoDB database connection
func (d *Database) Close() error {
	if err := d.Client.Disconnect(context.Background()); err != nil {
		zap.L().Error("failed to disconnect from MongoDB", zap.Error(err))
		return err
	}
	zap.L().Info("disconnected from MongoDB")
	return nil
}
