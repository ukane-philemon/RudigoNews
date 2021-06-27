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
//change this to var pagectex = context.TODO() for local development.
var pagectex context.Context

func init() {
	pagectex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
  defer cancel()
      //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(pagectex, clientOptions)
   client, err := mongo.Connect(pagectex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(pagectex, nil)
	if err != nil {
		log.Println(err)
	}

	pagecoll = client.Database("pageDB").Collection("pages")
}

type Page struct {
	ID              primitive.ObjectID
	Title           string
	Slug            string
	RawContent      template.HTML
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

func UpdatePage(pageId primitive.ObjectID, page Page) error {

	filter := bson.D{primitive.E{Key: "id", Value: pageId}}
	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "slug", Value: page.Slug},
		primitive.E{Key: "title", Value: page.Title},
		primitive.E{Key: "rawcontent", Value: page.RawContent},
		primitive.E{Key: "pagedescription", Value: page.PageDescription},
		primitive.E{Key: "author", Value:page.Author},
		primitive.E{Key: "datepublished", Value: page.DatePublished},
		primitive.E{Key: "datemodified", Value: page.DateModified},
		primitive.E{Key: "tags", Value: page.Tags},
		primitive.E{Key: "readtime", Value: page.ReadTime},
	}}}

	_, err := pagecoll.UpdateOne(pagectex, filter, updater)
	if err != nil {
		return err
	}

	return nil
}


func DeletePage (categoryId primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "id", Value: categoryId}}
	_ = pagecoll.FindOneAndDelete(pagectex, filter)
	return nil
}

func GetPage(pageslug string) (page Page, err error) {
	page = Page{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "slug", Value: pageslug}}

	err = pagecoll.FindOne(pagectex, filter).Decode(&page)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	//Return result without any error.
	return page, err
}
