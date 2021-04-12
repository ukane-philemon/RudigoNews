package main

import (
	//"context"
	"fmt"
	//"log"
	"net/http"

	//"github.com/gorilla/mux"

	admin "github.com/ukane-philemon/RudigoNews/admin_controllers"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"

	rout "github.com/ukane-philemon/RudigoNews/blogpost"
	//_ "github.com/ukane-philemon/RudigoNews/models"
)
/* 
var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("adminDB").Collection("User")
} */
func main() {
	mux := http.NewServeMux()

	mux.Handle("/upload/", http.StripPrefix("/upload/", http.FileServer(http.Dir("upload"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/admin-assets/", http.StripPrefix("/admin-assets/", http.FileServer(http.Dir("admin-assets"))))

	//Blog Routes

	mux.HandleFunc("/", rout.Home)
	mux.HandleFunc("/category/news", rout.News)
	mux.HandleFunc("/category/jobs", rout.Jobs)
	mux.HandleFunc("/category/tech", rout.Tech)
	mux.HandleFunc("/privacy-policy", rout.Privacy)
	//mux.HandleFunc("/search", rout.Search)
	mux.HandleFunc("/login", rout.Login)
	mux.HandleFunc("/about", rout.About)
	mux.HandleFunc("/contact", rout.ContactUs)
	mux.HandleFunc("/disclaimer", rout.Disclaimer)
	mux.HandleFunc("/post/{category}/{id:[0-9]+}", rout.PageHandler)
	mux.HandleFunc("/password-reset", rout.PasswordReset)
	//Admin Routesfr

	mux.HandleFunc("/admin/dashboard", admin.Dashboard)
	mux.HandleFunc("/logout", admin.Logout)
	mux.HandleFunc("/admin/all-post", admin.Allpost)
	mux.HandleFunc("/admin/new-post", admin.Newpost)
	mux.HandleFunc("/admin/comments", admin.Comments)
	mux.HandleFunc("/admin/profile", admin.Profile)
	mux.HandleFunc("/admin/update-post", admin.Updatepost)
	fmt.Println("Started on port 3000")
	http.ListenAndServe(":3000", mux)

	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// r := mux.NewRouter()
	// r.HandleFunc("/", rout.Home)
	// r.HandleFunc("/about", rout.About)

	//http.ListenAndServe(":3000", mux)
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
