package blogpost

import (
	"html/template"
	"net/http"
)

func Disclaimer(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/disclaimer.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
