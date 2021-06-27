package admin_controllers

import (
	"fmt"
	"html/template"
	"log"

	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"strings"
	"time"

	readingtime "github.com/ukane-philemon/RudigoNews/utils/readingtime"
	"github.com/ukane-philemon/RudigoNews/utils/saveimage"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Newpost handle displaying newpost editor and accepts post requests for sending data to db.
func Newpost(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		//get current user from session then find in db.
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if user.LoginState {
			posts := model.GetPosts()
			categories := model.GetCategories()
			type detail struct {
				Profile    model.User
				Articles   []model.Post
				Categories []model.Category
			}

			details := detail{
				Profile:    user,
				Articles:   posts,
				Categories: categories,
			}

			if err != nil {
				fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
				return
			}

			tmp, err := template.ParseFiles(
				"admin/template/editortemplate.gohtml",
				"admin/template/sidebar.gohtml",
				"admin/template/header.gohtml",
				"admin/template/editorfooter.gohtml",
				"admin/newpost.gohtml",
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
		//get current user from session then .find in db
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)

		if err != nil {
			//handle err
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		//check if user is loggedin before performing create new post request.
		if user.LoginState {

			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

			request.ParseMultipartForm(100000)

			//get values from form.
			Title := request.FormValue("title")
			Rawcontent := request.FormValue("content")
			Author := request.FormValue("author")
			Category := request.FormValue("category")
			Tags := request.FormValue("tags")
			Slug := request.FormValue("slug")
			PostDescription := request.FormValue("postdescription")
			FeaturedImage, header, err := request.FormFile("featuredImage")

			if err != nil {
				//handle err
				fmt.Fprintln(response, err, "No image found")
				return
			}

			//estimating read time
			estimation := readingtime.Estimate(Rawcontent)
			ReadTime := estimation.Text
			//save featured image to the upload folder
			im, err := saveimage.SaveImage(FeaturedImage, header)
			if err != nil {
				fmt.Fprintln(response, err)
				return
			}

			RawContent := template.HTML(Rawcontent)
			Slug = strings.ToLower(strings.Replace(Slug, " ", "-", -1))
			//check if post counter is in db if not, create it.
			count, err := model.GetCounterByName("Posts")
			if err != nil || count == nil {
				model.CreateCount(model.Counter{
					Name:   "Posts",
					Number: 0,
				})
			}
			//feed values to struct
			post := model.Post{
				ID:              primitive.NewObjectID(),
				Title:           Title,
				Slug:            Slug,
				RawContent:      RawContent,
				Category:        Category,
				Author:          Author,
				FeaturedImage:   header.Filename,
				ImageWidth:      im.Width,
				ImageHeight:     im.Height,
				PostDescription: PostDescription,
				ReadTime:        ReadTime,
				Tags:            Tags,
				DatePublished:   time.Now(),
				DateModified:    time.Now(),
			}

			postcount, err := model.GetCounterByName("Posts")

			if err != nil {
				log.Println(err)
			}
			//create post.
			err = model.CreatePost(post)

			if err != nil {
				//handle err
				fmt.Fprint(response, err)
			} else {
				//increase post count.
				model.AddCount(postcount.Name, postcount.Number+1)
			}

			http.Redirect(response, request, "/admin/new-post", http.StatusFound)

		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
		}
	default:
		http.Redirect(response, request, "/login", http.StatusForbidden)
	}
}
