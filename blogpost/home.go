package blogpost

import (
	"html/template"
	"net/http"
)

func Home(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/index.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
