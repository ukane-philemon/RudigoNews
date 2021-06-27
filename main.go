package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	admin "github.com/ukane-philemon/RudigoNews/admin_controllers"
	route "github.com/ukane-philemon/RudigoNews/blogpost"
)

func main() {
	mux := mux.NewRouter()

	mux.PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("upload"))))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.PathPrefix("/admin-assets/").Handler(http.StripPrefix("/admin-assets/", http.FileServer(http.Dir("admin-assets"))))

	//Admin Routes

	mux.HandleFunc("/admin/dashboard", admin.Dashboard)
	mux.HandleFunc("/logout", admin.Logout)
	mux.HandleFunc("/admin/all-post", admin.Allpost)
	mux.HandleFunc("/admin/all-pages", admin.Allpages)
	mux.HandleFunc("/admin/all-users", admin.AllUsers)
	mux.HandleFunc("/admin/new-post", admin.Newpost)
	mux.HandleFunc("/admin/comments", admin.Comments)
	mux.HandleFunc("/admin/categories", admin.Category)
	mux.HandleFunc("/admin/createpage", admin.Createpage)
	mux.HandleFunc("/admin/profile", admin.Profile)
	mux.HandleFunc("/admin/update", admin.Update)
	mux.HandleFunc("/admin/delete", admin.Delete)
	mux.HandleFunc("/admin/searchresult", admin.AdminSearch)
	mux.HandleFunc("/admin/task", admin.Task)
	//Blog Routes

	mux.HandleFunc("/", route.Home)
	mux.HandleFunc("/category/{category}", route.CategoryHandler)
	mux.HandleFunc("/searchresult", route.Search)
	mux.HandleFunc("/password-reset", route.PasswordReset)
	mux.HandleFunc("/login", route.Login)
	mux.HandleFunc("/newuser", route.Newuser)
	mux.HandleFunc("/contact", route.ContactUs)
	mux.HandleFunc("/subscribe", route.Subscribe)
	mux.HandleFunc("/{page}", route.PageHandler)
	mux.HandleFunc("/{category}/{slug}", route.PostPageHandler)

	fmt.Println("Started on port 3000")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, mux)

}

// /* var t *template.Template
// var routeMatch *regexp.Regexp
// var pd pageData */
// func servePage(w http.ResponseWriter, r *http.Request) {
//  /*  matches := routeMatch.FindStringSubmatch(r.URL.Path)

//   if len(matches) >= 1 {
//     page := matches[1] + ".html"
//     if t.Lookup(page) != nil {
//       w.WriteHeader(200)
//       t.ExecuteTemplate(w, page, pd)
//       return
//     }
//   } else if r.URL.Path == "/" {
//     w.WriteHeader(200)
//     t.ExecuteTemplate(w, "index.html", pd)
//     return
//   }
//   w.WriteHeader(404)
//   w.Write([]byte("NOT FOUND")) */
// }
//
