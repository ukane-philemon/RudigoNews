package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	model "github.com/ukane-philemon/RudigoNews/models"
)

type Detail struct {
	Page       model.Page
	Articles   []model.Post
	Loggedin   bool
	Categories []model.Category
	Pages      []model.Page
}

//Pagehandler handles page activities.
func PageHandler(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":

		vars := mux.Vars(request)
		slug := vars["page"]
		page, pageerr := model.GetPage(slug)

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

		Details := Detail{
			Page:       page,
			Articles:   posts,
			Loggedin:   user.LoginState,
			Categories: categories,
			Pages:      Pages,
		}

		if pageerr != nil {
			tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
				"view/404.html",
				"view/template/footer.html",
				"view/template/header.html",
				"view/template/sidebar.html",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "layout", Details)
			return
		}

		tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
			"view/template/pagetemplate.html",
			"view/template/footer.html",
			"view/template/header.html",
			"view/template/sidebar.html",
		)

		if err != nil {

			log.Println(err)
			http.Error(response, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmp.ExecuteTemplate(response, "layout", Details)

		return

	case "POST":

		fmt.Fprint(response, "Cannot send Post Resquest")
		http.Redirect(response, request, "/", http.StatusBadRequest)
		return

	default:

		fmt.Fprint(response, "Method not allowed")
		http.Redirect(response, request, "/", http.StatusBadRequest)
		return

	}

}
