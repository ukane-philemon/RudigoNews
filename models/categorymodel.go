package model

import (
	"context"
	"html/template"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var categorycoll *mongo.Collection
//change this to var categoryctex = context.TODO() for local development.
var categoryctex context.Context

func init() {
	categoryctex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
  defer cancel()
      //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(categoryctex, clientOptions)
  client, err := mongo.Connect(categoryctex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(categoryctex, nil)
	if err != nil {
		log.Println(err)
	}
	

	categorycoll = client.Database("pageDB").Collection("categories")

	
}

type Category struct {
	ID                  primitive.ObjectID
	Name                string
	Slug                string
	Author              string
	CategoryDescription template.HTML
	DatePublished       time.Time
	DateModified        time.Time
}
//CreateCategory creates a new category in categories db.
func CreateCategory(category Category) error {

	_, err := categorycoll.InsertOne(categoryctex, category)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

//GetCategories returns all categories in collection in an array.
func GetCategories() []Category {

opts := options.Find().SetSort(bson.M{"$natural": -1})
	cursor, err := categorycoll.Find(categoryctex, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var category []Category
	if err = cursor.All(categoryctex, &category); err != nil {
		log.Fatal(err)
	}
	return category

}

//Get CaegoryBrySlug gets a single category from the categories collection.
func GetCategoryBySlug(categoryslug string) (category *Category, err error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "slug", Value: categoryslug}}

	err = categorycoll.FindOne(categoryctex, filter).Decode(&category)

	//Return result without any error.
	return category, err
}

//GetCategoryByName returns any category whose name value is passed in as argument value.
func GetCategoryByName(categoryname string) (category *Category, err error) {
	//Define filter query for fetching specific document from collection
	
	filter := bson.D{primitive.E{Key: "name", Value: categoryname}}

	err = categorycoll.FindOne(categoryctex, filter).Decode(&category)

	//Return result without any error.
	return category, err
}

//UpdateCategory updates any category by Id.
func UpdateCategory(categoryId primitive.ObjectID, category Category) error {

	filter := bson.D{primitive.E{Key: "id", Value: categoryId}}

	_, err := categorycoll.ReplaceOne(categoryctex, filter, category)
	if err != nil {
		return err
	}

	return nil
}

//Deletes category by Id.
func DeleteCategory(categoryId primitive.ObjectID) (removedCategory Category, err error) {
	filter := bson.D{primitive.E{Key: "id", Value: categoryId}}
	err = categorycoll.FindOneAndDelete(categoryctex, filter).Decode(&removedCategory)
	return removedCategory, err
}
