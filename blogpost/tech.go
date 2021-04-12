package blogpost

import (
	"html/template"
	"net/http"
)

func Tech(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/tech.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}