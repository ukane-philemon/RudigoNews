package blogpost

import (
	"html/template"
	"log"
	"net/http"

	"github.com/hanzoai/gochimp3"
)

const (
	apiKey = "219d92820495a61675e232dd05b995e3-us19"
)

//subscribe handles subscribtion to mailchip.
func Subscribe(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":

		client := gochimp3.New(apiKey)

		// Audience ID
		listID := "b74452e89d"

		// Fetch list
		list, err := client.GetList(listID, nil)
		if err != nil {
			log.Printf("Failed to get list %s, %s", listID, err)

		}

		request.ParseForm()
		email := request.FormValue("email")
		// Add subscriber
		if email != "" {
			req := &gochimp3.MemberRequest{
				EmailAddress: email,
				Status:       "subscribed",
			}

			if _, err := list.CreateMember(req); err != nil {
				tmp, err := template.ParseFiles(
					"view/suberror.gohtml",
				)

				if err != nil {
					log.Print(err)
					http.Error(response, "internal server error", http.StatusInternalServerError)
					return
				}

				tmp.Execute(response, nil)

			} else {
				tmp, err := template.ParseFiles(
					"view/subsuccess.gohtml",
				)
				if err != nil {
					log.Print(err)
					http.Error(response, "internal server error", http.StatusInternalServerError)
					return
				}
				tmp.Execute(response, req.EmailAddress)
				return

			}
		} else {
			http.Redirect(response, request, "/", http.StatusFound)
		}

	default:
		http.Redirect(response, request, "/", http.StatusFound)
	}

}
