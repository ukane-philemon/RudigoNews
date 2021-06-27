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
	readingtime "github.com/ukane-philemon/RudigoNews/utils/readingtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Createpage accepts get request for displaying creating new page editor and post request for sending the  new post to db
func Createpage(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		pages := model.Getpages()
		type detail struct {
			Profile model.User
			Pages   []model.Page
		}

		details := detail{
			Profile: user,
			Pages:   pages,
		}

		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {
			tmp, err := template.ParseFiles(
				"admin/template/editortemplate.html",
				"admin/template/sidebar.html",
				"admin/template/header.html",
				"admin/template/editorfooter.html",
				"admin/createpage.html",
			)

			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}

			tmp.ExecuteTemplate(response, "layout", details)

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
		}

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {
			count, err := model.GetCounterByName("Pages")
			if err != nil || count == nil {
				model.CreateCount(model.Counter{
					Name:   "Pages",
					Number: 0,
				})
			}

			request.ParseForm()

			Title := strings.TrimSpace(request.FormValue("title"))
			Rawcontent := request.FormValue("content")
			Author := request.FormValue("author")
			DatePublished := time.Now()
			Tags := request.FormValue("tags")
			Slug := request.FormValue("slug")
			PageDescription := strings.TrimSpace(request.FormValue("pagedescription"))

			//estimating read time
			estimation := readingtime.Estimate(Rawcontent)
			fmt.Println(estimation.Text)     // "1 min read"
			fmt.Println(estimation.Duration) // 1 min
			fmt.Println(estimation.Words)    // amount (500)
			ReadTime := estimation.Text

			RawContent := template.HTML(Rawcontent)
			Slug = strings.ToLower(strings.Replace(Slug, " ", "-", -1))
			page := model.Page{
				ID:              primitive.NewObjectID(),
				Title:           Title,
				Slug:            Slug,
				RawContent:      RawContent,
				Author:          Author,
				PageDescription: PageDescription,
				ReadTime:        ReadTime,
				Tags:            Tags,
				DatePublished:   DatePublished,
				DateModified:    time.Now(),
			}
			pagecount, err := model.GetCounterByName("Pages")
			if err != nil {
				log.Println(err)
			}

			err = model.CreatePage(page)
			if err != nil {
				fmt.Fprint(response, err)
			} else {
				model.AddCount(pagecount.Name, pagecount.Number+1)
			}
			http.Redirect(response, request, "/admin/createpage", http.StatusFound)

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
		}
	default:
		http.Redirect(response, request, "/admin/dashboard", http.StatusFound)

	}

}
