package blogpost

import (
	"fmt"
	"html/template"
	"net/http"

	model "github.com/ukane-philemon/RudigoNews/models"
)


func PasswordReset(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
	tmp, _ := template.ParseFiles(
		"view/passwordreset.html",
	)

	tmp.ExecuteTemplate(response, "passwordreset.html", nil)

case "POST": 
err := request.ParseForm()
		if err != nil {
			fmt.Print(err) // Handle error here via logging and then return
		}

		userName := request.FormValue("username")
		PassWord := request.FormValue("password")
		
		model.PasswordHash(PassWord)
		if 	model.UpdateUserPassword(userName, PassWord) == nil {
				http.Redirect(response, request, "/admin/dashboard", http.StatusSeeOther)
		}
		
	}
}