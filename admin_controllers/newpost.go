package admin_controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Newpost(response http.ResponseWriter, request *http.Request) {

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

		// if err != nil {
		// 	fmt.Fprint(response, err)
		// }
		// if user.LoginState {

		tmp, _ := template.ParseFiles(
			"admin/newpost.html",
		)
		tmp.ExecuteTemplate(response, "newpost.html", details)
		// } else {
		// 	http.Redirect(response, request, "/login", http.StatusForbidden)
		// }

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		request.ParseForm()
		title := request.FormValue("title")
		Rawcontent := request.FormValue("content")
		Author := request.FormValue("author")
		Category := request.FormValue("category")
		date := time.Now()

		fmt.Printf(
			`title:%s
			Rawcontent:%s
			Author:%s
			Category:%s
			date:=%s
			`, title, Rawcontent, Author, Category, date,
		)
		fmt.Print(date)

		post := model.Post{
			Id:         primitive.ObjectID{},
			Title:      title,
			RawContent: Rawcontent,
			Category:   Category,
			Author:     Author,
			Date:       date,
		}

		model.CreatePost(post)
		http.Redirect(response, request, "/admin/new-post", http.StatusFound)
	}
}
