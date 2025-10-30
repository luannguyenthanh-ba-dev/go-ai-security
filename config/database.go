package config

// Database configuration

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDatabase creates a new MongoDB database connection
func NewMongoDatabase(uri, database string) (*Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		zap.L().Error("failed to connect to MongoDB", zap.String("uri", uri), zap.String("database", database), zap.Error(err))
		return nil, err
	}

	// Ping the database to check if the connection is successful
	if err := client.Ping(context.Background(), nil); err != nil {
		zap.L().Error("failed to ping MongoDB", zap.String("uri", uri), zap.String("database", database), zap.Error(err))
		return nil, err
	}

	return &Database{
		Client:   client,
		Database: client.Database(database),
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
