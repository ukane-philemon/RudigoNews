package model

import (
	"context"
	"log"
	"time"

	//	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var countercoll *mongo.Collection
//change this to var counterctex = context.TODO() for local development.
var counterctex = context.TODO()

func init() {
	counterctex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	    //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(counterctex, clientOptions)
	 client, err := mongo.Connect(counterctex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(counterctex, nil)
	if err != nil {
		log.Println(err)
	}

	countercoll = client.Database("pageDB").Collection("Counter")
}

type Counter struct {
	Name   string
	Number int
}

func CreateCount (count Counter) error {
	_, err := countercoll.InsertOne(counterctex, count)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func GetCounterByName(name string) (count *Counter, err error) {
	
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	err = countercoll.FindOne(counterctex, filter).Decode(&count)

	//Return result without any error.
	return count, err
}

func GetCounters() []Counter {

	cursor, err := countercoll.Find(counterctex, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var counts []Counter
	if err = cursor.All(counterctex, &counts); err != nil {
		log.Fatal(err)
	}
	return counts

}

func DeleteCount(name string, value int) error {
	//Define filter query for fetching specific document from countercoll
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "number", Value: value},
	}}}

	//Perform UpdateOne operation & validate against the error.
	_, err := countercoll.UpdateOne(counterctex, filter, updater)
	if err != nil {
		return err
	}
	return err
}

func AddCount(name string, value int) error {
	//Define filter query for fetching specific document from countercoll
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "number", Value: value},
	}}}

	//Perform UpdateOne operation & validate against the error.
	_, err := countercoll.UpdateOne(counterctex, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return err
}
