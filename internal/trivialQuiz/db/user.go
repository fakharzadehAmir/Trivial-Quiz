package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Birthday time.Time          `bson:"birthday"`
	Email    string             `bson:"email"`
}

// String returns a string representing the fields of User
func (u *User) String() string {
	return fmt.Sprintf("User{%+v}", *u)
}

// CreateNewUser creates a new user in the database MongoDB
func (mdb *MongoDB) CreateNewUser(ctx *context.Context, newUser *User) error {
	//	Check if user with given username exists
	existedUser := mdb.Collections.
		UserCollection.Collection.
		FindOne(*ctx, bson.M{"username": newUser.Username})

	if existedUser.Err() != nil {
		if errors.Is(existedUser.Err(), mongo.ErrNoDocuments) {
			// Encrypting the user password
			if encryptedPW, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4); err != nil {
				return err
			} else {
				newUser.Password = string(encryptedPW)
			}
			//	User with the same username doesn't exist, proceed to create a new user
			_, err := mdb.Collections.
				UserCollection.Collection.InsertOne(*ctx, newUser)
			if err != nil {
				return err
			}
			//	User created successfully
			return nil
		}
		//	Handle other errors that might have occurred during query execution
		return existedUser.Err()
	}

	//	User with the same username already exists
	return errors.New("this username already exists")
}

func (mdb *MongoDB) GetUserByUsername(ctx *context.Context, username string) (*User, error) {

	//	Check existence of user with given username
	findUser := mdb.Collections.
		UserCollection.Collection.
		FindOne(*ctx, bson.M{"username": username})

	//	Handle error of retrieving the user with given username
	if findUser.Err() != nil {
		return nil, findUser.Err()
	}

	//	Decode existed user to the declared variable existedUser
	var existedUser = &User{}
	if err := findUser.Decode(existedUser); err != nil {
		return nil, err
	}

	return existedUser, nil
}
