package blogpost

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ContactUs handles contact us page display and post request sends comment to db.
func ContactUs(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		posts := model.GetPosts()
		userName := GetUserName(request)
		user, err := model.GetUser(userName)
		categories := model.GetCategories()
		Pages := model.Getpages()

		if err != nil {
			user = model.User{}
		}

		Details := Detail{
			Page:       model.Page{},
			Articles:   posts,
			Loggedin:   user.LoginState,
			Categories: categories,
			Pages:      Pages,
		}

		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
			"slice" : func (array []interface{}, start int, end int) []interface{} {
			sliced := array[start:end]
			return sliced
			},
		}

		tmp, err := template.New(" ").Funcs(funcMap).ParseFiles(
			"view/template/footer.gohtml",
			"view/template/header.gohtml",
			"view/contact.gohtml",
		)

		if err != nil {
			//handle err
			log.Println(err)
			http.Error(response, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmp.ExecuteTemplate(response, "layout", Details)
		return

	case "POST":

		request.ParseForm()
		//get values from form and save to varibles.
		name := request.FormValue("Name")
		subject := request.FormValue("Subject")
		message := request.FormValue("Message")
		email := request.FormValue("Email")
		date := time.Now()

		//check if comment counter is availble, if not create it.

		count, err := model.GetCounterByName("Comments")
		if err != nil || count == nil {
			model.CreateCount(model.Counter{
				Name:   "Comments",
				Number: 0,
			})
		}

		comment := model.Comment{
			ID:           primitive.NewObjectID(),
			Subject:      subject,
			Message:      message,
			Email:        email,
			Name:         name,
			DateReceived: date,
		}

		commentcount, err := model.GetCounterByName("Comments")

		if err != nil {
			//handle err
			log.Println(err)
		}

		err = model.CreateComment(comment)

		if err != nil {
			fmt.Fprint(response, err)
		} else {
			//increase comment count.
			model.AddCount(commentcount.Name, commentcount.Number+1)
		}

		http.Redirect(response, request, "/contact", http.StatusFound)
		return

	default:

		fmt.Fprint(response, "Method not allowed")

	}

}
