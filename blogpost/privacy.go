package blogpost

import (
	"html/template"
	"net/http"
)

func Privacy(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/template/template.html",
		"view/privacy.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}
