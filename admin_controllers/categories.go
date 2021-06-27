package admin_controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
//Category takes get request to display all categories and new category page  and post request to create a new category
func Category(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		//get current user from session
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		if user.LoginState {
			//uncategoriezed cannot be deleted, this code will check if uncategorized category is in db, if not, it will create it.
			uncategorized, err := model.GetCategoryByName("Uncategorized")
			if err != nil || uncategorized == nil {
				model.CreateCategory(model.Category{
					ID:                  primitive.NewObjectID(),
					Name:                "Uncategorized",
					Slug:                "uncategorized",
					Author:              "default",
					CategoryDescription: "Uncategorized category",
					DatePublished:       time.Now(),
					DateModified:        time.Now(),
				})
			}

			
			categories := model.GetCategories()
			type detail struct {
				Profile    model.User
				Categories []model.Category
			}

			details := detail{
				Profile:    user,
				Categories: categories,
			}

			tmp, err := template.ParseFiles(
				"admin/template/template.gohtml",
				"admin/template/sidebar.gohtml",
				"admin/template/header.gohtml",
				"admin/template/footer.gohtml",
				"admin/category.gohtml",
			)
			
			if err != nil {
				log.Println(err)
				http.Error(response, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmp.ExecuteTemplate(response, "layout", details)
			return

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
			return
		}

	case "POST":
		//get current username from session then find in db.
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

			request.ParseForm()

			Name := request.FormValue("name")
			Author := request.FormValue("author")
			DatePublished := time.Now()
			DateModified := time.Now()
			Slug := request.FormValue("slug")
			CategoryDescription := request.FormValue("categorydescription")

			Slug = strings.ToLower(strings.Replace(Slug, " ", "-", -1))
			Name = strings.Title(Name)
			Categorydescription := template.HTML(CategoryDescription)
			category := model.Category{
				ID:                  primitive.NewObjectID(),
				Name:                Name,
				Slug:                Slug,
				Author:              Author,
				CategoryDescription: Categorydescription,
				DatePublished:       DatePublished,
				DateModified:        DateModified,
			}
			//check if category counter is in db else create one 
			count, err := model.GetCounterByName("Categories")
			if err != nil || count == nil {
				model.CreateCount(model.Counter{
					Name:   "Categories",
					Number: 0,
				})
			}
			//check if category is already in db if an err is returned, it means the category is not in db then it is created but if no error is returned from db, post request is ignored.
			available, err := model.GetCategoryByName(category.Name)
			if err != nil || available == nil {
			categorycount, _ := model.GetCounterByName("Categories")
			err := model.CreateCategory(category)
			if err != nil {
				fmt.Fprint(response, err)
			} else {
				model.AddCount(categorycount.Name, categorycount.Number+1)
			}
				
				
				if err != nil {
					fmt.Fprint(response, err)
				}
			}

			http.Redirect(response, request, "/admin/categories", http.StatusFound)
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
