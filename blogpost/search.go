package blogpost

import (
	"html/template"
	"net/http"
)

func Search(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/search.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
