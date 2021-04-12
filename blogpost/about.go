package blogpost

import (
	"net/http"
	"html/template"
)


func About(response http.ResponseWriter, request *http.Request) {
	tmp, _:= template.ParseFiles (
		"view/template/template.html",
		"view/about.html",
	)
	
	tmp.ExecuteTemplate(response, "layout", nil)
}