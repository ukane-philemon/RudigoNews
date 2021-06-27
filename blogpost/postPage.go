package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	model "github.com/ukane-philemon/RudigoNews/models"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type pdetail struct {
	Post       model.Post
	Articles   []model.Post
	Loggedin   bool
	Categories []model.Category
	Pages      []model.Page
}

//Postpagehandler takes care of the post to be display on request.
func PostPageHandler(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":

		vars := mux.Vars(request)
		postslug := string(vars["slug"])
		post, perr := model.GetPost(postslug)
		posts := model.GetPosts()
		userName := GetUserName(request)
		user, err := model.GetUser(userName)
		categories := model.GetCategories()
		Pages := model.Getpages()

		if err != nil {
			user = model.User{}
		}

		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
			"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
		}

		detail := pdetail{
			Post:       post,
			Articles:   posts,
			Loggedin:   user.LoginState,
			Categories: categories,
			Pages:      Pages,
		}

		if perr != nil {

			tmp, _ := template.New(" ").Funcs(funcMap).ParseFiles(
				"view/404.gohtml",
				"view/template/footer.gohtml",
				"view/template/header.gohtml",
				"view/template/sidebar.gohtml")
			tmp.ExecuteTemplate(response, "layout", detail)
			return

		}
		tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
			"view/post.gohtml",
			"view/template/footer.gohtml",
			"view/template/header.gohtml",
			"view/template/sidebar.gohtml",
		)

		if err != nil {
			log.Print(err)
			http.Error(response, "internal server error", http.StatusInternalServerError)
			return
		}

		tmp.ExecuteTemplate(response, "layout", detail)
		model.AddPostCount(post.Slug, post.Views+1)
		return

	case "POST":

		fmt.Fprint(response, "Cannot send Post Resquest")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	default:
		fmt.Fprint(response, "Method not allowed")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	}
}
