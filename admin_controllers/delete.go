package admin_controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ukane-philemon/RudigoNews/blogpost"
	model "github.com/ukane-philemon/RudigoNews/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Delete function handles all delete action performed by an admin, it does not permit delete requested by a non-admin user.
//Delete accepts only post requests.
func Delete(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	//Code for POST requests
	case "POST":
		//Get username details from blogpost module then use it to get user details from DB
		userName := blogpost.GetUserName(request)
		user, err := model.GetUser(userName)
		//check if any error in getting user details
		if err != nil {
			//respond with link to login if username is not in blogpost module memory
			fmt.Fprint(response, `<p>You are not currently Logged in, Click <a href="/login">Here</a> to Log in </p>`)
			return
		}
		//check if user state is true (logged in)
		if user.LoginState {
			//next the form is parse and relevant information are extracted
			request.ParseForm()
			postId, _ := primitive.ObjectIDFromHex(request.FormValue("postId"))
			action := request.FormValue("action")
			//check the action type
			switch action {
			case "categorydelete":
				if user.Adminrights {
					removedCategory, err := model.DeleteCategory(postId)

					if err != nil {
						log.Println(err)
						return
					} else {
						categorycount, _ := model.GetCounterByName("Categories")
						model.DeleteCount(categorycount.Name, categorycount.Number-1)
						model.ChangeCategorytoDefault(removedCategory)
						http.Redirect(response, request, "/admin/categories", http.StatusFound)
					}
				}

			case "postdelete":
				if user.Adminrights {
					post, err := model.DeletePost(postId)
					if err == nil {
					os.Remove("upload/" + post.FeaturedImage)
					postcount, _ := model.GetCounterByName("Posts")
					if postcount.Number <= 0 {
						postcount.Number = 0
						model.DeleteCount(postcount.Name, postcount.Number)
					} else {
						model.DeleteCount(postcount.Name, postcount.Number-1)
					}
				
					}
				}
				http.Redirect(response, request, "/admin/all-post", http.StatusFound)

			case "taskdelete":
				if user.Adminrights {
					model.DeleteTask(postId)
					taskcount, err := model.GetCounterByName("Tasks")
					if err == nil {
					model.DeleteCount(taskcount.Name, taskcount.Number-1)
					}
				}

				http.Redirect(response, request, "/admin/dashboard", http.StatusFound)

			case "commentdelete":
				if user.Adminrights {
					model.DeleteComment(postId)
					commentcount, _ := model.GetCounterByName("Comments")
					model.DeleteCount(commentcount.Name, commentcount.Number-1)
				}

				http.Redirect(response, request, "/admin/comments", http.StatusFound)

			case "pagedelete":
				if user.Adminrights {
					model.DeletePage(postId)
					pagescount, _ := model.GetCounterByName("Pages")
					model.DeleteCount(pagescount.Name, pagescount.Number-1)
				}

				http.Redirect(response, request, "/admin/all-pages", http.StatusFound)
			case "userdelete":
				if user.Adminrights {
					user, _ := model.DeleteUser(postId)
					if user.Avatar != "default.png" {
						os.Remove("upload/" + user.Avatar)
					}
					usercount, _ := model.GetCounterByName("Users")
					model.DeleteCount(usercount.Name, usercount.Number-1)
				}
				http.Redirect(response, request, "/admin/all-users", http.StatusFound)
			default:
				fmt.Fprint(response, "method not allowed")
			}
			return

		} else {
			//If user is not logged in
			http.Redirect(response, request, "/login", http.StatusSeeOther)
		}

	default:
		//if it is not a post request
		http.Redirect(response, request, "/", http.StatusSeeOther)

	}
}
