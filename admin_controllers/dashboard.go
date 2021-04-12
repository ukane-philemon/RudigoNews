package admin_controllers

import (
	//"fmt"
	"fmt"
	"html/template"
	"net/http"

	//	"github.com/ukane-philemon/RudigoNews/blogpost"
	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

func Dashboard(response http.ResponseWriter, request *http.Request) {

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
			"admin/template/template.html",
			"admin/dashboard.html",
		)
		tmp.ExecuteTemplate(response, "layout", details)
		// } else {
		// 	http.Redirect(response, request, "/", http.StatusForbidden)
		// }
	}

}
