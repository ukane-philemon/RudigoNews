package blogpost

import (
	"html/template"
	"net/http"
	)

type Page struct {
	Title   string
	Content string
	Date    string
	Id 		int
}
 
// use post struct in getting value from db, then pass it as option to template loader then utilize in frontend
func PageHandler(response http.ResponseWriter, request *http.Request) {
	tmp, err := template.ParseFiles(
		"view/template/template.html",
		"view/post.html",
	)
if err != nil {
		fileName := "view/404.html"
		http.ServeFile(response, request, fileName)
	}

tmp.ExecuteTemplate(response, "layout", nil)

}

