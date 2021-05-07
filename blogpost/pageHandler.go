package blogpost

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/ukane-philemon/RudigoNews/models"
)

func PageHandler(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		vars := mux.Vars(request)
		slug := vars["page"]
		page, err := model.Getpage(slug)
			if err != nil {
			tmp, _ := template.ParseFiles("view/404.html")
			tmp.Execute(response, nil)
		} else {			
		tmp, _ := template.ParseFiles("view/template/pagetemplate.html")
	
			tmp.Execute(response, page)
		}

	case "POST":
		fmt.Fprint(response, "Cannot send Post Resquest")
		http.Redirect(response, request, "/", http.StatusBadRequest)
	default:
		fmt.Fprint(response, "Method not allowed")
		http.Redirect(response, request, "/", http.StatusBadRequest)

	}

}