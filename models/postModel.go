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

var coll *mongo.Collection
//change this to var ctex = context.TODO() for local development.
var ctex context.Context

func init() {
	ctex, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	    //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(ctex, clientOptions)
	 client, err := mongo.Connect(ctex, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(ctex, nil)
	if err != nil {
		log.Println(err)
	}

	coll = client.Database("PostDB").Collection("Posts")
	
	//creating index for search
	opt := options.Index()
	opt.SetUnique(false)
	opt.SetName("postsearch")
	opt.SetWeights(bson.M{
		"slug":            2,
		"featuredimage":   1,
		"title":           5, // Word matches in the title are weighted 5× standard.
		"rawcontent":      4, // Word matches in the title are weighted 4× standard.
		"postdescription": 2,
	})

	index := mongo.IndexModel{Keys: bson.D{
		primitive.E{Key: "slug", Value: "text"},
		primitive.E{Key: "title", Value: "text"},
		primitive.E{Key: "featuredimage", Value: "text"},
		primitive.E{Key: "rawcontent", Value: "text"},
		primitive.E{Key: "author", Value: "text"},
		primitive.E{Key: "datepublished", Value: "text"},
		primitive.E{Key: "datemodified", Value: "text"},
		primitive.E{Key: "tags", Value: "text"},
		primitive.E{Key: "views", Value: "text"},
		primitive.E{Key: "category", Value: "text"},
		primitive.E{Key: "postdescription", Value: "text"},
	}, Options: opt}

	if _, err := coll.Indexes().CreateOne(ctex, index); err != nil {
		log.Println("Could not create index:", err)
	} else {
		log.Println("success")
	}

}

type Post struct {
	ID              primitive.ObjectID `bson:"id"`
	Title           string             `bson:"title"`
	Views           int                `bson:"views"`
	Slug            string             `bson:"slug"`
	RawContent      template.HTML      `bson:"rawcontent"`
	Category        string             `bson:"category"`
	Author          string             `bson:"author"`
	FeaturedImage   string             `bson:"featuredimage"`
	ImageWidth      int                `bson:"imagewidth"`
	ImageHeight     int                `bson:"imageheight"`
	PostDescription string             `bson:"postdescription"`
	Tags            string             `bson:"tags"`
	ReadTime        string             `bson:"readtime"`
	DatePublished   time.Time          `bson:"datepubished"`
	DateModified    time.Time          `bson:"datemodified"`
}

func CreatePost(post Post) error {

	_, err := coll.InsertOne(ctex, post)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil

}

func GetPosts() []Post {
	opts := options.Find().SetSort(bson.M{"$natural": -1})
	cursor, err := coll.Find(ctex, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var posts []Post
	if err = cursor.All(ctex, &posts); err != nil {
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

func DeletePost(postId primitive.ObjectID) (post Post, err error) {

	//Define filter query for deleting specific document from collection
	filter := bson.D{primitive.E{Key: "id", Value: postId}}

	err = coll.FindOneAndDelete(ctex, filter).Decode(&post)

	return post, err
}

func UpdatePost(postId primitive.ObjectID, post Post) error {

	filter := bson.D{primitive.E{Key: "id", Value: postId}}
	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "category", Value: post.Category},
		primitive.E{Key: "slug", Value: post.Slug},
		primitive.E{Key: "title", Value: post.Title},
		primitive.E{Key: "featuredimage", Value: post.FeaturedImage},
		primitive.E{Key: "rawcontent", Value: post.RawContent},
		primitive.E{Key: "author", Value: post.Author},
		primitive.E{Key: "datepublished", Value: post.DatePublished},
		primitive.E{Key: "datemodified", Value: post.DateModified},
		primitive.E{Key: "tags", Value: post.Tags},
		primitive.E{Key: "imageheight", Value: post.ImageHeight},
		primitive.E{Key: "imagewidth", Value: post.ImageWidth},
		primitive.E{Key: "readtime", Value: post.ReadTime},
	}}}

	_, err := coll.UpdateOne(ctex, filter, updater)
	if err != nil {
		return err
	}

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

func GetPost(postslug string) (post Post, err error) {

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "slug", Value: postslug}}

	err = coll.FindOne(ctex, filter).Decode(&post)

	//Return result without any error.
	return post, err
}

func ChangeCategorytoDefault(removedCategory Category) error {
	filter := bson.D{primitive.E{Key: "category", Value: removedCategory.Name}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "category", Value: "Uncategorized"},
	}}}
	//Perform UpdateOne operation & validate against the error.
	_, err := coll.UpdateMany(ctex, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func AddPostCount(slug string, value int) error {
	//Define filter query for fetching specific document from countercoll
	filter := bson.D{primitive.E{Key: "slug", Value: slug}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "views", Value: value},
	}}}

	//Perform UpdateOne operation & validate against the error.
	_, err := coll.UpdateOne(ctex, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return err
}

func TextSearch(searchterm string) ([]Post, error) {
	filter := bson.M{"$text": bson.M{"$search": searchterm}}
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"slug":            2,
		"featuredimage":   1,
		"title":           5, // Word matches in the title are weighted 5× standard.
		"rawcontent":      4, // Word matches in the title are weighted 5× standard.
		"postdescription": 2,
		"author":          1,
		"datepublished":   1,
		"datemodified":    1,
		"category":        1,
		"score":           bson.M{"$meta": "textScore"},
	})
	findOptions.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := coll.Find(ctex, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var posts []Post
	if err = cursor.All(ctex, &posts); err != nil {
		log.Fatal(err)
	}
	return posts, err
	
}
