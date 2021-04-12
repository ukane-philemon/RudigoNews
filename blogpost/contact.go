package blogpost

import (
	"html/template"
	"net/http"
)

func ContactUs(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles(
		"view/contact.html",
	)

	tmp.ExecuteTemplate(response, "contact.html", nil)
}
