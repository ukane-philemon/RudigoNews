package admin_controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

//Comments is responsible for getting comments from db and diplaying to a loggedin user. It accepts only get requests.
func Comments(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		//get current user from session the find in db.
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		//check if user is loggedin
		if user.LoginState {
			comments := model.GetComments()
			type detail struct {
				Profile  model.User
				Comments []model.Comment
			}

			details := detail{
				Profile:  user,
				Comments: comments,
			}

			if err != nil {
				fmt.Fprint(response, err)
			}

			tmp, err := template.ParseFiles(
				"admin/template/template.html",
				"admin/template/sidebar.html",
				"admin/template/header.html",
				"admin/template/footer.html",
				"admin/comments.html",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "layout", details)

		} else {
			http.Redirect(response, request, "/", http.StatusForbidden)
		}
	default:
		fmt.Fprint(response, "Not allowed")
	}

}
