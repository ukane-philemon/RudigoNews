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

var commentcoll *mongo.Collection
//change this to var commentctex = context.TODO() for local development.
var commentctex = context.TODO()

func init() {
	commentctex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
  defer cancel()
      //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(commentctex, clientOptions)
   client, err := mongo.Connect(commentctex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(commentctex, nil)
	if err != nil {
		log.Println(err)
	}

	commentcoll = client.Database("pageDB").Collection("comments")
}

type Comment struct {
	
	ID              primitive.ObjectID
	Subject         string
	Message         string
	Email			string
	Name        	string
	DateReceived    time.Time
	
}

//CreateComment adds a new comment to commentcoll.
func CreateComment(comment Comment) error {

	_, err := commentcoll.InsertOne(commentctex, comment)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

//GetComments return an array of all comments in commentcoll.
func GetComments() []Comment {

	cursor, err := commentcoll.Find(commentctex, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var comments []Comment
	if err = cursor.All(commentctex, &comments); err != nil {
		log.Fatal(err)
	}
	return comments

}


//DeleteComment Removes one comment from commentcoll using comment Id
func DeleteComment(Id primitive.ObjectID) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "id", Value: Id}}
_ = commentcoll.FindOneAndDelete(commentctex, filter)
	return nil

}
