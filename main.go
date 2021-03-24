package main

import (
	"net/http"
	"github.com/ukane-philemon/RudigoNews/routes/homeroute"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeroute.Index)
	http.HandleFunc("/home", homeroute.Index)
	http.HandleFunc("/home/index", homeroute.Index)

	http.ListenAndServe(":3000", nil)
}

// /* var all = template.Must(template.ParseFiles(
// 	"index1.html",
// 	"login.html",
// 	"disclaimer.html",
// 	"contact.html",
// 	"news.html",
// 	"education.html",
// 	"jobs.html",
// 	"privacy.html",

// 	))

// /* type data struct {
// 	Result string
// }

// var result string */

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "index1.html", nil)
// }

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "news.html", nil)
// }

// func educationHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "education.html", nil)
// }

// func jobsHandler(w http.ResponseWriter, r *http.Request) {
// all.ExecuteTemplate(w, "jobs.html", nil)
// }

// func privacyHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "privacy.html", nil)
// }

// func disclaimerHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "disclaimer.html", nil)

// }

// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	all.ExecuteTemplate(w, "about.html", nil)
// }

// package main

// import (
//   "html/template"
//   "log"
//   "net/http"
//   "os"
//   "regexp"
// )

// /* type pageData struct {
//   Title       string
//   CompanyName string
// } */

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
// func main() {
//   var err error
//   t, err = template.ParseGlob("./public/*")
//   if err != nil {
//     log.Println("Cannot parse templates:", err)
//     os.Exit(-1)
//   }
//   /* routeMatch, _ = regexp.Compile(`^\/(\w+)`)
//   pd = pageData{
//     "WidgetCo Home",
//     "WidgetCo International",
//   } */
//   http.HandleFunc("/", servePage)
//   log.Fatal(http.ListenAndServe(":8080", nil))
// }
