package admin_controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Profile(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		user, err := model.GetUser(blogpost.Username)
		if err != nil {
			fmt.Fprint(response, err)
		}

		posts := model.GetPost()
		type detail struct {
			Profile  model.User
			Articles []model.Post
		}

		details := detail{
			Profile:  user,
			Articles: posts,
		}

		if user.LoginState {
			tmp, _ := template.ParseFiles(
				"admin/template/template.html",
				"admin/profile.html",
			)

			tmp.ExecuteTemplate(response, "layout", details)
			// } else {
			// 	http.Redirect(response, request, "/", http.StatusForbidden)
		}

	case "POST":
		request.ParseMultipartForm(100000)

		userName := request.FormValue("userName")
		firstname := request.FormValue("firstName")
		lastName := request.FormValue("lastName")
		PassWord := request.FormValue("password")
		Email := request.FormValue("email")
		address := request.FormValue("address")
		City := request.FormValue("city")
		State := request.FormValue("state")
		Zip := request.FormValue("zip")
		pics, header, err := request.FormFile("pics")

		if err != nil {
			fmt.Fprintln(response, err)
			return
		}
		defer pics.Close()
		filetype := header.Header.Get("Content-Type")

		switch filetype {
		case "image/jpeg", "image/jpg", "image/gif", "image/png":
			// Creating file in folder
			out, err := os.Create("upload/" + header.Filename)
			if err != nil {
				fmt.Fprintf(response, "Unable to create the file for writing. Check your write access privilege")
				return
			}
			defer out.Close()
			// write the content from POST REQUEST to the file
			_, err = io.Copy(out, pics)
			if err != nil {
				fmt.Fprintln(response, err)
			}
			uploadpath := "upload/" + header.Filename
			PassWord = model.PasswordHash(PassWord)

			newUser := model.User{
				ID:       primitive.NewObjectID(),
				UserName: userName,
				Email:    Email,
				First:    firstname,
				Last:     lastName,
				Password: PassWord,
				Address:  address,
				Avatar:   uploadpath,
				City:     City,
				State:    State,
				Zip:      Zip,
			}

			if model.UpdateUser(newUser) == nil {
				http.Redirect(response, request, "/admin/dashboard", http.StatusSeeOther)
				fmt.Print("user updated")

			} else {
				model.CreateUser(newUser)
			}

		default:
			fmt.Fprint(response, "Please upload an image")

		}

	}

}
