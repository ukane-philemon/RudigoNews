package admin_controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

func Deletepost(response http.ResponseWriter, request *http.Request) {

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
		tmp, _ := template.ParseFiles(
			"admin/newpost.html",
		)
		tmp.ExecuteTemplate(response, "newpost.html", details)

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		title := request.FormValue("title")
		content := request.FormValue("content")
		fmt.Printf("title = %s\n", title)
		fmt.Printf("content = %s\n", content)
		http.Redirect(response, request, "/admin/new-post", http.StatusFound)
	}
}
