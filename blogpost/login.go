package blogpost

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"

	model "github.com/ukane-philemon/RudigoNews/models"
	"golang.org/x/crypto/bcrypt"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:       "session",
			Value:      encoded,
			Path:       "/",
			Domain:     "",
			Expires:    time.Now().Add(6000 * time.Second),
			MaxAge:     60 * 60 * 2,
			Secure:     false,
			HttpOnly:   false,
			
		}
		http.SetCookie(response, cookie)
	}
}

func ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// Login handler takes care of user loggin and setting of cookie

func Login(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	//Code for GET requests
	case "GET":

		userName := GetUserName(request)
		user, err := model.GetUser(userName)

		if userName != "" && user.LoginState {
			http.Redirect(response, request, "/admin/dashboard", http.StatusFound)
			return
		} else if err != nil || !user.LoginState {
			tmp, err := template.ParseFiles(
				"view/login.html",
			)
			if err != nil {
				log.Print(err)
				http.Error(response, "internal server error", http.StatusInternalServerError)
				return
			}
			tmp.ExecuteTemplate(response, "login.html", nil)
			return
		}

	//Code for POST requests
	case "POST":

		request.ParseForm()

		Username := request.FormValue("username")
		password := request.FormValue("password")
		if Username != "" && password != "" {
			// .. checking credentials ..
			//Check if user is available and get user from DB
			user, err1 := model.GetUser(Username)

			//convert saved and the given password to byte so bcrytp can use it
			hashedPassword := []byte(user.Password)
			Password := []byte(password)
			//compare
			err2 := bcrypt.CompareHashAndPassword(hashedPassword, Password)
			//give err if not match
			if err1 != nil || err2 != nil {

				error := "Username or Password is Incorrect, Try again."
				tmp, _ := template.ParseFiles(
					"view/login.html",
				)
				tmp.ExecuteTemplate(response, "login.html", error)
				return
			} else {
				//set login state
				model.LoginState(Username, true)
				setSession(Username, response)

				//give access if match
				http.Redirect(response, request, "/admin/dashboard", http.StatusSeeOther)
				//fmt.Print(user.Password)
			}
		}
	}
}
