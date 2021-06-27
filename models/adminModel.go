package model

import (
	"time"

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
//change this to var ctx = context.TODO() for local development.
var ctx context.Context

func init() {
	
   ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
  defer cancel()
    //use this for local developement 
//  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
// 	client, err := mongo.Connect(ctx, clientOptions)
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(
     "mongodb+srv://<cluster name>:<password>@<cluster link>/adminDB?retryWrites=true&w=majority",
  ))

	if err != nil {
		log.Println(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	collection = client.Database("adminDB").Collection("users")
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
	Adminrights bool
	DateJoined  time.Time
}

//CreateUser adds a new user to db.
func CreateUser(user User) error {

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//GetUser gets one particular user from db using username filter.
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

//Updateuser updates users profile information.
func UpdateUser(user User, ID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "id", Value:ID}}

	_, err := collection.ReplaceOne(ctx, filter, user)
	if err != nil {
		return err
	}

	return nil
}

//UpdateUserPassword changes/updates user password.
func UpdateUserPassword(username, password string) (*mongo.UpdateResult, error) {

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "password", Value: password},
	}}}

	//Perform UpdateOne operation & validate against the error.
	result, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return result, err
	}
	//Return success without any error.
	return result, err
}

//LoginState sets and unsets user login state.
func LoginState(username string, loginState bool) error {

	filter := bson.D{primitive.E{Key: "username", Value: username}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "loginstate", Value: loginState},
	}}}
	//Perform UpdateOne operation & validate against the error.
	_, err := collection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//MakeAdmin makes a user an admin, Only admins can delete items from blog.
func MakeAdmin(userId primitive.ObjectID) error {
filter := bson.D{primitive.E{Key: "id", Value: userId}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "adminrights", Value: true},
	}}}
	//Perform UpdateOne operation & validate against the error.
	_, err := collection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//RemoveAdmin removes a user from admin status.
 func RemoveAdmin(userId primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "id", Value: userId}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "adminrights", Value: false},
	}}}
	//Perform UpdateOne operation & validate against the error.
	_, err := collection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
 }

 //PasswordHash hashs user password and returns a bcrypt form of it.
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

//GetUsers gets all users from db and returns an array of results.
func GetUsers() []User {
	opts := options.Find().SetSort(bson.M{"$natural": -1})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users

}

//DelteUser Removes user from db.
func DeleteUser (userId primitive.ObjectID) (user User, err error) {
	filter := bson.D{primitive.E{Key: "id", Value: userId}}
	err = collection.FindOneAndDelete(ctx, filter).Decode(&user)
	return user, err
}

