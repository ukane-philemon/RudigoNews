package homeroute

import (
	"html/template"
	"net/http"
)

func Index(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"../public/view/index.html",
		"../public/template/template.html",
	)

	tmp.ExecuteTemplate(response, "layout", nil)
}