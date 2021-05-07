package blogpost

import (
	"fmt"
	"html/template"
	"net/http"

	model "github.com/ukane-philemon/RudigoNews/models"
)

func Business(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		posts := model.GetPost()
		var bposts []model.Post
		for _, p := range posts {
			if p.Category == "Business" {
				bposts = append(bposts, p)

			}
		}
		tmp, _ := template.ParseFiles(
			"view/template/template.html",
			"view/business.html",
		)

		tmp.ExecuteTemplate(response, "layout", bposts)

	case "POST":
		fmt.Fprint(response, "Cannot Send Post Request")
		http.Redirect(response, request, "/", http.StatusNotAcceptable)

	default:
		fmt.Fprint(response, "Method Not Allowed")
		http.Redirect(response, request, "/", http.StatusForbidden)

	}
}
