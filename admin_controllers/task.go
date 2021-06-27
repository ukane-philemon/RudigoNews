package admin_controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Task(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		if user.LoginState {
			count, err := model.GetCounterByName("Tasks")
			if err != nil || count == nil {
				model.CreateCount(model.Counter{
					Name:   "Tasks",
					Number: 0,
				})
			}
			request.ParseForm()

			Title := request.FormValue("title")
			Author := request.FormValue("author")
			DatePublished := time.Now()
			TaskDescription := request.FormValue("taskdescription")

			
			task := model.Task{
				ID:            primitive.NewObjectID(),
				Title:         Title,
				Description:   TaskDescription,
				Author:        Author,
				DatePublished: DatePublished,
			}
			taskcount, err := model.GetCounterByName("Tasks")

			if err != nil {
				log.Println(err)
			}

			err = model.CreateTask(task)
			
			if err != nil {
				fmt.Fprint(response, err)
			} else {
				model.AddCount(taskcount.Name, taskcount.Number+1)
			}

			

			http.Redirect(response, request, "/admin/dashboard", http.StatusFound)
			return

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
			return
		}
	default:
		http.Redirect(response, request, "/admin/dashboard", http.StatusFound)
		return
	}

}
