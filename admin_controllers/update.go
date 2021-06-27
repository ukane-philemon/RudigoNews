package admin_controllers

import (
	"fmt"
	"html/template"
	"image"
	"log"
	"net/http"
	"strings"
	"time"

	//"github.com/ukane-philemon/RudigoNews/blogpost"
	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	readingtime "github.com/ukane-philemon/RudigoNews/utils/readingtime"
	save "github.com/ukane-philemon/RudigoNews/utils/saveimage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Update handles all updates/edit requests.
func Update(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}

		if user.LoginState {
			postslug := request.URL.Query().Get("postslug")
			action := request.URL.Query().Get("action")
			switch action {
			case "postedit":
				post, err := model.GetPost(postslug)
				categories := model.GetCategories()
				if err != nil {
					tmp, _ := template.ParseFiles("view/404.html")
					tmp.Execute(response, nil)
					return
				}
				type detail struct {
					Profile    model.User
					Post       model.Post
					Categories []model.Category
				}

				details := detail{
					Profile:    user,
					Post:       post,
					Categories: categories,
				}

				tmp, err := template.ParseFiles(
					"admin/template/editortemplate.html",
					"admin/template/sidebar.html",
					"admin/template/header.html",
					"admin/template/editorfooter.html",
					"admin/updatepost.html",
				)

				if err != nil {
					log.Println(err)
					fmt.Fprintln(response, err)
					return
				}

				tmp.ExecuteTemplate(response, "layout", details)
				return

			case "categoryedit":
				category, err := model.GetCategoryBySlug(postslug)
				if err != nil {
					tmp, _ := template.ParseFiles("view/404.html")
					tmp.Execute(response, nil)
					return
				}
				type detail struct {
					Profile  model.User
					Category model.Category
				}

				details := &detail{
					Profile:  user,
					Category: *category,
				}

				tmp, _ := template.ParseFiles(
					"admin/template/editortemplate.html",
					"admin/template/sidebar.html",
					"admin/template/header.html",
					"admin/template/editorfooter.html",
					"admin/updatecategory.html",
				)
				tmp.ExecuteTemplate(response, "layout", details)
				return

			case "pageedit":
				page, err := model.GetPage(postslug)

				if err != nil {
					tmp, _ := template.ParseFiles("view/404.html")
					tmp.Execute(response, nil)
					return
				}

				type detail struct {
					Profile model.User
					Page    model.Page
				}

				details := &detail{
					Profile: user,
					Page:    page,
				}

				tmp, err := template.ParseFiles(
					"admin/template/editortemplate.html",
					"admin/template/sidebar.html",
					"admin/template/header.html",
					"admin/template/editorfooter.html",
					"admin/updatepage.html",
				)

				if err != nil {
					log.Print(err)
					http.Error(response, "internal server error", http.StatusInternalServerError)
					return
				}

				tmp.ExecuteTemplate(response, "layout", details)
				return

			default:
				fmt.Fprintln(response, "Not found and Not recognized")
			}
		} else {

			http.Redirect(response, request, "/login", http.StatusContinue)

		}

	case "POST":
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		if err != nil {
			fmt.Fprint(response, `<p>You are not Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		if user.LoginState {
			request.ParseMultipartForm(100000)
			updatetype := request.FormValue("updatetype")
			switch updatetype {
			case "postupdate":
				Title := request.FormValue("title")
				postID, _ := primitive.ObjectIDFromHex(request.FormValue("postId"))
				Rawcontent := request.FormValue("content")
				Author := request.FormValue("author")
				Category := request.FormValue("category")
				DatePublished, err := time.Parse("Jan 02, 2006 3:04 PM", request.FormValue("datepublished"))
				if err != nil {
					fmt.Println(err)
				}

				DateModified := time.Now()
				Tags := request.FormValue("tags")
				Slug := request.FormValue("slug")
				PostDescription := request.FormValue("postdescription")
				FeaturedImage, header, err := request.FormFile("featuredImage")

				var im image.Config
				var featuredimage string
				if err != nil {
					featuredimage = request.FormValue("featuredimage")
				} else {
					im, _ = save.SaveImage(FeaturedImage, header)
					featuredimage = header.Filename
					if err != nil {
						log.Println(err)

					}
				}

				RawContent := template.HTML(Rawcontent)
				Slug = strings.ToLower(strings.Replace(Slug, " ", "-", -1))
				//estimating read time
				estimation := readingtime.Estimate(Rawcontent)
				ReadTime := estimation.Text

				post := model.Post{
					ID:              postID,
					Title:           Title,
					Slug:            Slug,
					RawContent:      RawContent,
					Category:        Category,
					Author:          Author,
					FeaturedImage:   featuredimage,
					ImageWidth:      im.Width,
					ImageHeight:     im.Height,
					PostDescription: PostDescription,
					ReadTime:        ReadTime,
					Tags:            Tags,
					DatePublished:   DatePublished,
					DateModified:    DateModified,
				}

				model.UpdatePost(postID, post)

				http.Redirect(response, request, "/admin/all-post", http.StatusFound)

				return
			case "pageupdate":
				request.ParseMultipartForm(100000)

				Title := request.FormValue("title")
				pageID, _ := primitive.ObjectIDFromHex(request.FormValue("pageId"))
				Rawcontent := request.FormValue("content")
				Author := request.FormValue("author")
				DatePublished, err := time.Parse("Jan 02, 2006 3:04 PM", request.FormValue("datepublished"))
				if err != nil {
					fmt.Println(err)
				}

				DateModified := time.Now()
				Tags := request.FormValue("tags")
				Slug := request.FormValue("slug")
				PageDescription := request.FormValue("pagedescription")

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
					Tags:            Tags,
					ReadTime:        ReadTime,
					DatePublished:   DatePublished,
					DateModified:    DateModified,
				}
				updaterr := model.UpdatePage(pageID, page)
				if err != nil {
					fmt.Fprint(response, updaterr)
				} else {
					http.Redirect(response, request, "/admin/dashboard", http.StatusFound)
				}
			case "categoryupdate":
				request.ParseForm()

				Name := strings.Title(request.FormValue("name"))
				Author := request.FormValue("author")
				DatePublished, _ := time.Parse("Jan 02, 2006 at 3:04PM", request.FormValue("datepublished"))
				CategoryId, _ := primitive.ObjectIDFromHex(request.FormValue("categoryId"))
				DateModified := time.Now()
				Slug := strings.ToLower(strings.Replace(request.FormValue("slug"), " ", "-", -1))
				CategoryDescription := strings.TrimSpace(request.FormValue("categorydescription"))

				Categorydescription := template.HTML(CategoryDescription)
				category := model.Category{
					ID:                  CategoryId,
					Name:                Name,
					Slug:                Slug,
					Author:              Author,
					CategoryDescription: Categorydescription,
					DatePublished:       DatePublished,
					DateModified:        DateModified,
				}

				err := model.UpdateCategory(CategoryId, category)
				if err != nil {
					fmt.Fprint(response, err)
				}

				http.Redirect(response, request, "/admin/categories", http.StatusFound)
				return

			default:
				fmt.Fprint(response, "Update Not Accepted")

			}
		} else {
			http.Redirect(response, request, "/login", http.StatusFound)
		}
	default:
		http.Redirect(response, request, "/login", http.StatusForbidden)
	}
}
