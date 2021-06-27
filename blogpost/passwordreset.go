package blogpost

import (
	"html/template"
	"log"
	"net/http"

	model "github.com/ukane-philemon/RudigoNews/models"
	"golang.org/x/crypto/bcrypt"
)

//PasswordReset as the name implies resets the user password.
func PasswordReset(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":

		tmp, _ := template.ParseFiles(
			"view/passwordreset.gohtml",
		)

		tmp.ExecuteTemplate(response, "passwordreset.gohtml", nil)

	case "POST":

		err := request.ParseForm()
		if err != nil {
			// Handle error here via logging and then return
			log.Println(err)
		}

		userName := request.FormValue("username")
		oldpassword := request.FormValue("oldpassword")
		PassWord := model.PasswordHash(request.FormValue("password"))

		user, err := model.GetUser(userName)
		

		if err != nil {

			errtext := "Username not found, Try again."
			tmp, err := template.ParseFiles(
				"view/passwordreset.gohtml",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "passwordreset.gohtml", errtext)
			return

		} 
		//convert saved and the given password to byte so bcrytp can use it
			hashedPassword := []byte(user.Password)
			oldPassword := []byte(oldpassword)
			//compare
		
		err = bcrypt.CompareHashAndPassword(hashedPassword, oldPassword)
		
		if err != nil {
			errtext := "Old Password not correct, Try again."
			tmp, err := template.ParseFiles(
				"view/passwordreset.gohtml",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "passwordreset.gohtml", errtext)
			return

		} else {

			result, _ := model.UpdateUserPassword(userName, PassWord)

			if result.MatchedCount == 0 {
				errtext := "Error, Please Try again."
				tmp, _ := template.ParseFiles(
					"view/passwordreset.gohtml",
				)
				tmp.ExecuteTemplate(response, "passwordreset.gohtml", errtext)
				return

			} else {

				http.Redirect(response, request, "/login", http.StatusSeeOther)
				return
			}
		}

	}
}
