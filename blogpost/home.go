package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	//"strings"

	model "github.com/ukane-philemon/RudigoNews/models"
)

type hdetail struct{
		Articles []model.Post
		Loggedin bool
		Categories []model.Category
		Pages   []model.Page

}

//Home handles all requests sent to the home page of the blog.

func Home(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		categories := model.GetCategories()
		posts := model.GetPosts()
		userName := GetUserName(request)
		user, err := model.GetUser(userName)
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

		hdetails := &hdetail{
			Articles:   posts,
			Loggedin:   user.LoginState,
			Categories: categories,
			Pages: Pages,
		}
			tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
				"view/template/template.gohtml",
				"view/template/footer.gohtml",
				"view/template/sidebar.gohtml",
				"view/template/header.gohtml",
				"view/index.gohtml",
			)

			if err != nil {
				log.Println(err)
				http.Error(response, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmp.ExecuteTemplate(response, "layout", hdetails)
		return

	case "POST":
		fmt.Fprint(response, "Cannot send Post Resquest")
		http.Redirect(response, request, "/", http.StatusBadRequest)
	default:
		fmt.Fprint(response, "Method not allowed")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	}

}
