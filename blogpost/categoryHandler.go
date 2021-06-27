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

type detail struct {
	Categoryposts []model.Post
	Articles      []model.Post
	Loggedin      bool
	Category      model.Category
	Categories    []model.Category
	Pages         []model.Page
}

//Categoryhandler takes care of all category activities for the blog.
func CategoryHandler(response http.ResponseWriter, request *http.Request) {

	switch request.Method {

	case "GET":

		vars := mux.Vars(request)
		categoryslug := string(vars["category"])
		category, categoryerr := model.GetCategoryBySlug(categoryslug)
		if categoryerr != nil {
			category = &model.Category{}
		}
		categories := model.GetCategories()
		Pages := model.Getpages()

		posts := model.GetPosts()
		userName := GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			user = model.User{}
		}

		var bposts []model.Post
		for _, p := range posts {
			if p.Category == category.Name {
				bposts = append(bposts, p)

			}
		}

		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
			"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
		}

		details := detail{
			Categoryposts: bposts,
			Articles:      posts,
			Loggedin:      user.LoginState,
			Category:      *category,
			Categories:    categories,
			Pages:         Pages,
		}

		if categoryerr != nil {
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
			tmp.ExecuteTemplate(response, "layout", details)
			return
		}
		tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
			"view/template/categorytemplate.html",
			"view/template/footer.html",
			"view/template/sidebar.html",
			"view/template/header.html",
		)
		if err != nil {
			//handle err
			log.Println(err)
			http.Error(response, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmp.ExecuteTemplate(response, "layout", details)

	case "POST":

		fmt.Fprint(response, "Cannot send Post Resquest")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	default:

		fmt.Fprint(response, "Method not allowed")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	}
}
