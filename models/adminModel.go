package model

import (
	_ "github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("adminDB").Collection("User")
}

type User struct {
	ID         primitive.ObjectID
	UserName   string
	Email      string
	First      string
	Last       string
	Password   string
	Address    string
	Avatar     string
	City       string
	State      string
	Zip        string
	LoginState bool
}

func CreateUser(user User) error {

	_, err := collection.InsertOne(ctx, user,  )
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func GetUser(username string) (User, error) {
	user := User{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	err := collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return user, err
	}
	//Return result without any error.
	return user, nil
}

func UpdateUser(user User) error {
	filter := bson.D{primitive.E{Key: "username", Value: user.UserName}}

	_, err := collection.ReplaceOne(ctx, filter, user)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserPassword(username, password string) error {

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "password", Value: PasswordHash(password)},
	}}}

	//Perform UpdateOne operation & validate against the error.
	_, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func LoginState(username string, loginState bool) error {

	filter := bson.D{primitive.E{Key: "username", Value: username}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "LoginState", Value: loginState},
	}}}
	//Perform UpdateOne operation & validate against the error.
	_, err := collection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func PasswordHash(pass string) string {
	password := []byte(pass)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	finalpass := string(hashedPassword)

	return finalpass
	// Comparing the password with the hash
	/* err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	   fmt.Println(err)  */ // nil means it is a match
}
