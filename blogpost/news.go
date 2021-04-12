package blogpost

import (
	"html/template"
	"net/http"
)

func News(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/news.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
