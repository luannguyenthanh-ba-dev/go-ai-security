package config

import "go.mongodb.org/mongo-driver/mongo"

type Database interface {
	Close() error
	GetMongoClient() *mongo.Client
	GetMongoDatabase() *mongo.Database
}
