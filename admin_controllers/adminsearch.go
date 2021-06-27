package admin_controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)

//AdminSearch function takes care of finding posts for a loggedin admin, it accepts only get requests.
func AdminSearch(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		UserName := blogpost.GetUserName(request)
		user, err := model.GetUser(UserName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {

			request.ParseForm()
			searchterm := request.FormValue("search")
			fmt.Println(searchterm)
			result, err := model.TextSearch(searchterm)
			if err != nil {
				fmt.Fprintln(response, err)
			}

			type detail struct {
				Profile  model.User
				Articles []model.Post
			}

			details := detail{
				Profile:  user,
				Articles: result,
			}

			tmp, err := template.ParseFiles(
				"admin/template/template.html",
				"admin/template/sidebar.html",
				"admin/template/header.html",
				"admin/template/footer.html",
				"admin/search.html",
			)

		if err != nil {
			log.Print(err)
			http.Error(response, "internal server error", http.StatusInternalServerError)
			return
		}

			tmp.ExecuteTemplate(response, "layout", details)
		}
	} else {
		fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
	
}
