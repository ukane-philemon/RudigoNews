package admin_controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"github.com/ukane-philemon/RudigoNews/utils/saveimage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Profile accepts post and get requests. It is responsible for displaying user profile details and post request for updating user details.
func Profile(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		type detail struct {
			Profile model.User
		}

		details := detail{
			Profile: user,
		}

		if user.LoginState {
			tmp, err := template.ParseFiles(
				"admin/template/template.gohtml",
				"admin/template/sidebar.gohtml",
				"admin/template/header.gohtml",
				"admin/template/footer.gohtml",
				"admin/profile.gohtml",
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

	case "POST":
		//get current user from session the nfind in db.
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			//handle err
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		//check if current user is loggedin.
		if user.LoginState {
			request.ParseMultipartForm(100000)

			userName := request.FormValue("userName")
			firstname := request.FormValue("firstName")
			lastName := request.FormValue("lastName")
			PassWord := request.FormValue("password")
			Email := request.FormValue("email")
			ID, _ := primitive.ObjectIDFromHex(request.FormValue("ID"))
			address := request.FormValue("address")
			City := request.FormValue("city")
			State := request.FormValue("state")
			Zip := request.FormValue("zip")
			pics, header, err := request.FormFile("pics")

			if err != nil {
				fmt.Fprintln(response, err)
				return
			}
			//save profile image to upload folder.
			//refer "github.com/ukane-philemon/RudigoNews/utils/saveimage"
			saveimage.SaveImage(pics, header)

			//Hash new password.
			PassWord = model.PasswordHash(PassWord)

			newUser := model.User{
				ID:       ID,
				UserName: userName,
				Email:    Email,
				First:    firstname,
				Last:     lastName,
				Password: PassWord,
				Address:  address,
				Avatar:   header.Filename,
				City:     City,
				State:    State,
				Zip:      Zip,
			}

			if model.UpdateUser(newUser, ID) == nil {
				http.Redirect(response, request, "/admin/profile", http.StatusSeeOther)
				fmt.Print("user updated")

			} else {
				model.CreateUser(newUser)
			}

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
		}
	default:
		http.Redirect(response, request, "/login", http.StatusForbidden)
	}

}
