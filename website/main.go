package robspychala

import (
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    template.HTML
	File    string
}

var indexTempl = template.Must(template.ParseFiles("static/index.html"))
var aboutTempl = template.Must(template.ParseFiles("static/about.html"))
var postTempl = template.Must(template.ParseFiles("static/post.html"))

func handleIndex(w http.ResponseWriter, r *http.Request) {
	posts := getPosts()
	err := indexTempl.Execute(w, posts)
	if err != nil {
		panic(err)
	}
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	posts := getPosts()
	err := aboutTempl.Execute(w, posts)
	if err != nil {
		panic(err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	f := "posts/" + r.URL.Path[1:] + ".md"
	fileread, _ := ioutil.ReadFile(f)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	date := string(lines[1])
	summary := string(lines[2])
	body := strings.Join(lines[3:len(lines)], "\n")
	body = string(blackfriday.MarkdownCommon([]byte(body)))
	post := Post{title, date, summary, template.HTML(body), r.URL.Path[1:]}
	err := postTempl.Execute(w, post)
	if err != nil {
		panic(err)
	}
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		a = append(a, Post{title, date, summary, template.HTML(body), file})
	}
	return a
}

func init() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", handleIndex).Methods("GET")
	rtr.HandleFunc("/about", handleAbout).Methods("GET")
	rtr.HandleFunc("/{name:[a-z]+}", handlePost).Methods("GET")

	http.Handle("/", rtr)
}
