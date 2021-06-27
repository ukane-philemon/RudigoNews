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

//Allpages takes care of getting all pages from the db and sending them to the admin dashborad in a table format.

func Allpages(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		//get current user from ssession
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		//check if user is logged in.
		if user.LoginState {
			//get pages from db if user is logged in.
			pages := model.Getpages()
			type detail struct {
				Profile model.User
				Pages   []model.Page
			}

			details := detail{
				Profile: model.User{},
				Pages:   pages,
			}

			if err != nil {
				fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
				return
			}

			funcMap := template.FuncMap{
				"ToLower": strings.ToLower,
				"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
			}

			tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
				"admin/template/template.html",
				"admin/template/sidebar.html",
				"admin/template/header.html",
				"admin/template/footer.html",
				"admin/viewallpages.html",
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
