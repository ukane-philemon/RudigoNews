package admin_controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

//Dashoard takes care of displaying data on the user dashboard.
func Dashboard(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		//get user from ssesion then find in db
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		posts := model.GetPosts()
		comments := model.GetComments()
		counters := model.GetCounters()
		tasks := model.GetTasks()

		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
			"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
		}
		type detail struct {
			Profile  model.User
			Articles []model.Post
			Comments []model.Comment
			Tasks    []model.Task
			Counters []model.Counter
		}

		details := &detail{
			Profile:  user,
			Articles: posts,
			Comments: comments,
			Tasks:    tasks,
			Counters: counters,
		}

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		//check if user is loggedin before displaying gohtml template.
		if user.LoginState {

			tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
				"admin/template/template.gohtml",
				"admin/template/sidebar.gohtml",
				"admin/template/header.gohtml",
				"admin/template/footer.gohtml",
				"admin/dashboard.gohtml",
			)

			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}

			tmp.ExecuteTemplate(response, "layout", details)
			return

		} else {

			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return

		}

	default:
		fmt.Fprint(response, "Not Allowed")
	}

}
