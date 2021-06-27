package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	model "github.com/ukane-philemon/RudigoNews/models"
)

type sdetail struct {
	Results    []model.Post
	Articles   []model.Post
	Loggedin   bool
	Categories []model.Category
	Pages      []model.Page
}

//Search handles all search request from the blog.
func Search(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":

		request.ParseForm()
		searchterm := request.FormValue("search")
		result, err := model.TextSearch(searchterm)
		Pages := model.Getpages()

		if err != nil {
			fmt.Fprintln(response, err)
		}

		categories := model.GetCategories()
		posts := model.GetPosts()
		userName := GetUserName(request)
		user, err := model.GetUser(userName)

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

		sdetails := &sdetail{
			Results:    result,
			Articles:   posts,
			Loggedin:   user.LoginState,
			Categories: categories,
			Pages:      Pages,
		}

		tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
			"view/template/template.html",
			"view/template/footer.html",
			"view/template/sidebar.html",
			"view/template/header.html",
			"view/search.html",
		)
		if err != nil {
			log.Print(err)
			http.Error(response, "internal server error", http.StatusInternalServerError)
			return
		}

		tmp.ExecuteTemplate(response, "layout", sdetails)
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

