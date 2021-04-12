package blogpost

import (
	"html/template"
	"net/http"
)

func Jobs(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/jobs.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
