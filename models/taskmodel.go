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

var taskcoll *mongo.Collection
//change this to var taskctex = context.TODO() for local development.
var taskctex context.Context

func init() {
	taskctex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
  defer cancel()
      //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(taskctex, clientOptions)
   client, err := mongo.Connect(taskctex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(taskctex, nil)
	if err != nil {
		log.Println(err)
	}

	taskcoll = client.Database("pageDB").Collection("tasks")
}

type Task struct {
	
	ID            primitive.ObjectID
	Title         string
	Description   string
	Author		  string
	DatePublished  time.Time
	
}

func CreateTask(task Task) error {

	_, err := taskcoll.InsertOne(taskctex, task)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

func GetTasks() []Task {

	cursor, err := taskcoll.Find(taskctex, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var tasks []Task
	if err = cursor.All(taskctex, &tasks); err != nil {
		log.Fatal(err)
	}
	return tasks

}



func DeleteTask(Id primitive.ObjectID) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "id", Value: Id}}
_ = taskcoll.FindOneAndDelete(taskctex, filter)
	return nil

}
