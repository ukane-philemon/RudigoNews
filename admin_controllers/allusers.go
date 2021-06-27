/*
The admin_controllers (as the name implies) are responsible for any action involving admin user action, admin dashboard actions.
*/

package admin_controllers

import (
	//	"fmt"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AllUsers gets all the users in the db and sends them in a struct to the user dashboard.
//It also receive posts requests to make a user admin or remove from admin.
func AllUsers(response http.ResponseWriter, request *http.Request) {
	//get current user from session then find in db.
	userName := blogpost.GetUserName(request)
	user, err := model.GetUser(userName)
	//check the method of request.
	switch request.Method {
	case "GET":
		users := model.GetUsers()
		type detail struct {
			Profile model.User
			Users   []model.User
		}

		details := detail{
			Profile: user,
			Users:   users,
		}

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {

			tmp, err := template.ParseFiles(
				"admin/template/template.html",
				"admin/template/sidebar.html",
				"admin/template/header.html",
				"admin/template/footer.html",
				"admin/viewallusers.html",
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
		request.ParseForm()
		userId, err := primitive.ObjectIDFromHex(request.FormValue("userId"))
		if err != nil {
			log.Println(err)
			http.Error(response, "User Id is not valid", http.StatusInternalServerError)
		}
		action := request.FormValue("action")

		if user.LoginState && user.Adminrights {
			switch action {
			case "remove":
				model.RemoveAdmin(userId)
				http.Redirect(response, request, "/admin/all-users", http.StatusFound)
			case "makeadmin":
				model.MakeAdmin(userId)
				http.Redirect(response, request, "/admin/all-users", http.StatusFound)
			}

		} else {
			http.Error(response, "You are not yet an Admin", http.StatusInternalServerError)
		}

	default:
		http.Redirect(response, request, "/admin/dashboard", http.StatusFound)

	}
}
