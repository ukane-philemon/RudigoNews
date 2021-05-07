package admin_controllers

import(
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	readingtime "github.com/ukane-philemon/RudigoNews/utils/readingtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Createpage(response http.ResponseWriter, request * http.Request) {
	switch request.Method {
		case "GET":
		user, err := model.GetUser(blogpost.Username)
		pages := model.Getpages()
		type detail struct {
			Profile model.User
			Pages[] model.Page
		}

		details	:= detail {
			Profile: user,
			Pages: pages,
		}

		if err != nil {
			fmt.Fprint(response, err)
			return
		}


		if user.LoginState {
			tmp,_:= template.ParseFiles(
				"admin/template/editortemplate.html",
				"admin/createpage.html",
			)
			tmp.ExecuteTemplate(response, "layout", details)

		} else {
	http.Redirect(response, request, "/login", http.StatusForbidden)
		}

case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		request.ParseForm()

		Title := request.FormValue("title")
		Rawcontent := request.FormValue("content")
		Author := request.FormValue("author")
		Category := request.FormValue("category")
		DatePublished := time.Now()
		Tags := request.FormValue("tags")
		Slug := request.FormValue("slug")
		PageDescription := request.FormValue("postdescription")



		//estimating read time
		estimation := readingtime.Estimate(Rawcontent)
		fmt.Println(estimation.Text) // "1 min read"
		fmt.Println(estimation.Duration) // 1 min
		fmt.Println(estimation.Words) // amount (500)
		ReadTime := estimation.Text


		RawContent := template.HTML(Rawcontent)
		Slug = strings.ToLower(strings.Replace(Slug, " ", "-", -1))
		page := model.Page {
			ID: primitive.NewObjectID(),
			Title: Title,
			Slug: Slug,
			RawContent: RawContent,
			Category: Category,
			Author: Author,
			PageDescription: PageDescription,
			ReadTime: ReadTime,
			Tags: Tags,
			DatePublished: DatePublished,
			DateModified: time.Time {},
		}

		err := model.CreatePage(page)
		if err != nil {
			fmt.Fprint(response,err)
		}
		http.Redirect(response, request, "/admin/createpage", http.StatusFound)
	
}
}