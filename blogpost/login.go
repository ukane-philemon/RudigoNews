package blogpost

import (
	"fmt"
	"html/template"
	"net/http"

	model "github.com/ukane-philemon/RudigoNews/models"
	"golang.org/x/crypto/bcrypt"
)

var Username string

func Login(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	//Code for GET requests
	case "GET":
		tmp, _ := template.ParseFiles(
			"view/login.html",
		)

		tmp.ExecuteTemplate(response, "login.html", nil)

	//Code for POST requests
	case "POST":

		err := request.ParseForm()
		if err != nil {
			fmt.Print(err)
		}

		Username = request.FormValue("username")
		password := request.FormValue("password")

		//Check if user is available and get user from DB
		user, err := model.GetUser(Username)

		if err != nil {
			fmt.Fprint(response, err)
		}
		//convert saved and the given password to byte so bcrytp can use it
		hashedPassword := []byte(user.Password)
		Password := []byte(password)
		//compare
		err2 := bcrypt.CompareHashAndPassword(hashedPassword, Password)
		//give err if not match
		if err2 != nil {
			fmt.Fprint(response, err2)
		} else {
			//set login state
			model.LoginState(Username, true)

			//give access if match
			http.Redirect(response, request, "/admin/dashboard", http.StatusSeeOther)
			//fmt.Print(user.Password)
		}
		fmt.Print (Username, password)
	}
}
