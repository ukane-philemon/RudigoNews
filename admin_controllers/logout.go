package admin_controllers

import (
	//"fmt"
	"net/http"

	// "github.com/ukane-philemon/RudigoNews/blogpost"
	// model "github.com/ukane-philemon/RudigoNews/models"
)



func Logout(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	//Code for GET requests
	case "GET":
		// user, err := model.GetUser(blogpost.Username)

		// if err != nil {
		// 	fmt.Fprint(response, err)
		// }
		// if user.LoginState {
		//model.LoginState(blogpost.Username, false)
		//}
			
		//http.Redirect(response, request, "/", http.StatusForbidden)
		
		
	//Code for POST requests
	case "POST":

	http.Redirect(response, request, "/", http.StatusForbidden)
		
}
}