package model

import (
	"context"
	"html/template"

	"log"
	"time"

	//	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

var pagecoll *mongo.Collection
var pagectex = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(pagectex, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(pagectex, nil)
	if err != nil {
		log.Fatal(err)
	}

	pagecoll = client.Database("pageDB").Collection("pages")
}

type Page struct {
	ID              primitive.ObjectID
	Title           string
	Slug            string
	RawContent      template.HTML
	Category        string
	Author          string
	PageDescription string
	Tags            string
	ReadTime		string
	DatePublished   time.Time
	DateModified    time.Time
}

func CreatePage(page Page) error {

	_, err := pagecoll.InsertOne(pagectex, page)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

func Getpages() []Page {

	cursor, err := pagecoll.Find(pagectex, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var page []Page
	if err = cursor.All(pagectex, &page); err != nil {
		log.Fatal(err)
	}
	return page

}

func UpdatePage(page Page) error {
	// title := post.Title

	// update := bson.D{

	// }

	// updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult)
	return nil
}


func Getpage(slug string) (page Page, err error) {
	page = Page{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "slug", Value: slug}}

	err = pagecoll.FindOne(pagectex, filter).Decode(&page)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	//Return result without any error.
	return page, err
}
