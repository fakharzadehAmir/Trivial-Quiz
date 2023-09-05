package db

import "go.mongodb.org/mongo-driver/mongo"

type Collections struct {
	UserCollection struct {
		Collection *mongo.Collection
		Name       string
	}
}
