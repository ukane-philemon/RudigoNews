package admin_controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

func Updatepost(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		user, err := model.GetUser(blogpost.Username)
		posts := model.GetPost()
		type detail struct {
			Profile  model.User
			Articles []model.Post
		}

		details := detail{
			Profile:  user,
			Articles: posts,
		}

		if err != nil {
			fmt.Fprint(response, err)
		}

		// if user.LoginState {
		tmp, _ := template.ParseFiles(
			"admin/updatepost.html",
		)
		tmp.ExecuteTemplate(response, "updatepost.html", details)

		// } else {
		// 	http.Redirect(response, request, "/login", http.StatusForbidden)
		// }

	}
}
