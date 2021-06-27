package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Newuser aevalutes and adds new users to db. It does not permit duplicate users username.
func Newuser(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":

		tmp, _ := template.ParseFiles(
			"view/newuser.html",
		)
		tmp.ExecuteTemplate(response, "newuser.html", nil)
		return

	case "POST":

		request.ParseMultipartForm(100000)
		//get values from form and save to varibles.
		userName := request.FormValue("username")
		firstname := request.FormValue("firstname")
		lastName := request.FormValue("lastname")
		PassWord := request.FormValue("password")
		Email := request.FormValue("email")

		//check if user exists, available is an error, it shows that the username does not exists therefore the account creation succeeds.
		_, available := model.GetUser(userName)

		if available != nil {
			//hash user password
			PassWord = model.PasswordHash(PassWord)

			newUser := model.User{
				ID:          primitive.NewObjectID(),
				UserName:    userName,
				Email:       Email,
				First:       firstname,
				Last:        lastName,
				Password:    PassWord,
				Avatar:      "default.png",
				Adminrights: false,
				DateJoined:  time.Now(),
			}

			usercount, err := model.GetCounterByName("Users")

			if err != nil {
				//handle err
				log.Println(err)
			}

			err = model.CreateUser(newUser)

			if err != nil {
				//handle err
				fmt.Fprint(response, err)
			} else {
				//increase users count.
				model.AddCount(usercount.Name, usercount.Number+1)
			}

			http.Redirect(response, request, "/login", http.StatusFound)
			return

		} else {

			usernametaken := "Username is not available, try something else :)"
			tmp, err := template.ParseFiles(
				"view/newuser.html",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "newuser.html", usernametaken)
			return

		}
	default:

		fmt.Println(response, "Method not accepted")
	}
}
