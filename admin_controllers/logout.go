package admin_controllers

import (
	"fmt"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
)
//Logout handles the logout event which involves setting user loginstate from db to false and deleting ssession cookie
func Logout(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	//Code for GET requests
	case "GET":
		//get user from session then find in db
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			fmt.Fprint(response, `<p>You are not currently Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		} else if user.LoginState {
			//this sets the current user loginstate to false
			model.LoginState(userName, false)
			fmt.Fprint(response, `<p>
			You have been logged out. If you wish to login again, Click <a href="/login">Here</a> to Log in.
			 </p>
			 <p>Click <a href="/">Here</a> to visit the blog.</p>
			`)
			//and clears session data
			blogpost.ClearSession(response)
			return

		} else {
			fmt.Fprint(response, `<p>You are not currently Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

	//Code for POST requests
	case "POST":

		http.Redirect(response, request, "/", http.StatusForbidden)

	default:

		http.Redirect(response, request, "/", http.StatusSeeOther)

	}
}
