package admin_controllers

import (
	//	"fmt"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

//Allpost takes care of getting all pages from the db and sending them to the admin dashborad in a table format. It accepts only get requests.
func Allpost(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		//get current user from session
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		//get posts from db
		posts := model.GetPosts()
		type detail struct {
			Profile  model.User
			Articles []model.Post
		}

		details := detail{
			Profile:  user,
			Articles: posts,
		}

		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
			"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
		}

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		//check if user is loggedin
		if user.LoginState {

			tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
				"admin/template/template.gohtml",
				"admin/template/sidebar.gohtml",
				"admin/template/header.gohtml",
				"admin/template/footer.gohtml",
				"admin/viewallpost.gohtml",
			)

			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}

			tmp.ExecuteTemplate(response, "layout", details)

		} else {
			http.Redirect(response, request, "/login", http.StatusSeeOther)
		}

	case "POST":
		fmt.Fprint(response, "Post Request Not Allowed")

	default:
		http.Redirect(response, request, "/admin/dashboard", http.StatusFound)

	}
}
