package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database wraps the MongoDB client and database
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect initializes the MongoDB connection
func Connect(uri, dbName string) (*Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(dbName)
	log.Println("Connected to MongoDB!")
	return &Database{
		Client:   client,
		Database: database,
	}, nil
}

// Disconnect closes the MongoDB connection
func (db *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("Disconnected from MongoDB!")

	return nil
}
