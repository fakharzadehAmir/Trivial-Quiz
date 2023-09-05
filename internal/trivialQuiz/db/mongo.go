package db

import (
	"Trivia_Quiz/config"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct {
	DB          *mongo.Database
	cfg         *config.Config
	logger      *logrus.Logger
	Collections *Collections
}

func ConnectDB(cfg *config.Config, logger *logrus.Logger) (*MongoDB, error) {
	// Create URL of MongoDB
	url := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/",
		cfg.Database.Username, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
	)
	// Set client options
	clientOptions := options.Client().ApplyURI(url)

	// Create a client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Connect to MongoDB container on docker
	err = client.Connect(ctx)
	if err != nil {
		defer client.Disconnect(ctx)
		return nil, err
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		defer client.Disconnect(ctx)
		return nil, err
	}

	return &MongoDB{
		DB:     client.Database(cfg.Database.DatabaseName),
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (mdb *MongoDB) CreateCollections() (error, *Collections) {

	var collections = &Collections{}
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//	Create "user" collection if it doesn't exist
	err, userCollection := assignCollection(ctx, mdb, "user")
	if err != nil {
		return err, nil
	}
	collections.UserCollection = struct {
		Collection *mongo.Collection
		Name       string
	}{Collection: userCollection, Name: "user"}
	// others

	return nil, collections
}

func assignCollection(ctx context.Context, mdb *MongoDB, collectionName string) (error, *mongo.Collection) {
	cursor, err := mdb.DB.ListCollections(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err, nil
	}
	defer cursor.Close(ctx)
	// Check if the collection already exists
	for cursor.Next(ctx) {
		var result struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err != nil {
			return err, nil
		}
		if result.Name == collectionName {
			mdb.logger.Infof("collection (\"%s\") exists in the database!", collectionName)
			return nil, mdb.DB.Collection(collectionName)
		}
	}

	// Create collection which is not in the database
	if err := mdb.DB.CreateCollection(ctx, collectionName); err != nil {
		return err, nil
	}
	mdb.logger.Infof("new collection (\"%s\") has been created!", collectionName)
	return nil, mdb.DB.Collection(collectionName)

}
