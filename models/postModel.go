package model

import (
	//"fmt"
	"net/http"
	"time"

	//"fmt"

	_ "github.com/gorilla/mux"

	"context"
	"log"

	//	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var ctex = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctex, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctex, nil)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database("postDB").Collection("posts")
}

type Post struct {
	Id         primitive.ObjectID
	Title      string
	RawContent string
	Category 	string
	Author		string
	Image		string
	Date       time.Time
}

func CreatePost(post Post) error {

	_, err := coll.InsertOne(ctex, post)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

func GetPost() []Post {

cursor, err := coll.Find(ctex, bson.D{})
if err != nil {
    log.Fatal(err)
}
var posts []Post
if err = cursor.All(ctx, &posts); err != nil {
    log.Fatal(err)
}
return posts

	// for cursor.Next(context.TODO()) {
	// 	Elem := &bson.D{}
	// 	if err := cursor.Decode(Elem); err != nil {
	// 		log.Fatal(err)
	// 	}
		
	// 	fmt.Println(Elem)
	

}

func UpdatePost(post Post) error {
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

/* func (p Post) TruncatedText() string {
	chars := 0
	for i, _ := range p.Content {
		chars++
		if chars > 150 {
			return p.Content[:i] + ` ...`
		}
	}
	return p.Content
} */

func ServePage(w http.ResponseWriter, r *http.Request) {

	// vars := mux.Vars(r)
	// pageGUID := vars["guid"]
	// thisPage := Page{}
	// fmt.Println(pageGUID)
	// err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	// thisPage.Content = template.HTML(thisPage.RawContent)
	// if err != nil {
	// 	http.Error(w, http.StatusText(404), http.StatusNotFound)
	// 	log.Println(err)
	// 	return
	// }
	// // html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`

	// t, _ := template.ParseFiles("templates/blog.html")
	// t.Execute(w, thisPage)
}
