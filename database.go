package main

import (
	"context"
	"os"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var database *mongo.Database

// getDatabase returns connected database. If none is connected, it connects to it first.
func getDatabase() *mongo.Database {
	if database != nil {
		return database
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGODB_USERNAME") + ":" + os.Getenv("MONGODB_PASSWORD") + "@cluster0.gkjrmx3.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil || client == nil {
		log.Fatal().Msg(err.Error())
	}
	return client.Database("db")
}

// readEvents tries to fetch events from the database
func readEvents() ([]Event, error) {
	db := getDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collectionExists(db, "event")
	if err != nil {
		return nil, err
	}

	cursor, err := db.Collection("event").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var events []Event
	err = cursor.All(ctx, &events)
	if err != nil {
		return nil, err
	}

	if len(events) != 0 {
		return events, nil
	}

	events = buildEvents()
	_, err = db.Collection("event").InsertMany(ctx, events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// collectionExists checks if a given collection exitst in a given database.
// If the collection does not exist, it is created.
func collectionExists(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get list of existing collections
	collections, err := db.ListCollectionNames(ctx, map[string]any{})
	if err != nil {
		return err
	}

	// Check if collection exists
	if slices.Contains(collections, collectionName) {
		return nil // Collection exists
	}

	// Create the collection
	err = db.CreateCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	return nil
}

func buildEvents() []Event {
	return []Event{
		{
			time.Date(2023, time.June, 30, 16, 45, 0, 0, time.Now().Location()),
			"Master degree",
			"Achieved master degree in IT-Management",
			Education,
		}, {
			time.Date(2020, time.September, 18, 17, 45, 0, 0, time.Now().Location()),
			"Bachelor degree",
			"Achieved bachelor degree in Business Informatics",
			Education,
		}, {
			time.Date(2017, time.August, 21, 8, 0, 0, 0, time.Now().Location()),
			"Dual Student @ HARTING Technology Group",
			"Started working as dual student Business Informatics within HARTING Technology Group. Mainly doing SAP development and PLM configuration.",
			Position,
		}, {
			time.Date(2017, time.July, 7, 0, 0, 0, 0, time.Now().Location()),
			"Finished school",
			"Finished school with advanced classes in physics and math",
			Education,
		},
	}
}
