package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/julienschmidt/httprouter"
)

// Model of stuff to render a page
type Model struct {
	Title string
	Name  string
}

// “export” means “public” => with upper case first letter
// the reason is simple: the renderer package use the reflect package in order to get/set fields values
// the reflect package can only access public/exported struct fields.
// So try defining
type pageData struct {
	Title  string
	Navbar string
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	data := pageData{Title: "Deon", Navbar: "hello"}
	tmpl := template.Must(template.ParseFiles("test.html"))
	err := tmpl.Execute(w, data)
	//err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func cst(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("cstest.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Header string
		Uname  string
	}{
		Title:  "Index Page",
		Header: "Hello, World!",
		Uname:  "Deon",
	}

	//Note use of ExecuteTemplate as apposed to just execute
	//err := tmpls.ExecuteTemplate(w, "index.html", data)
	//var tmpls = template.Must(template.ParseFiles("index.html", "headernav.html"))

	tmpl := template.Must(template.ParseFiles("index.html", "headernav.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("index executed")
}

func about(w http.ResponseWriter, r *http.Request) {
	s2 := template.Must(template.ParseFiles("about.html", "header.html", "headernav.html"))
	s2.Execute(w, nil)
	fmt.Println("about executed")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/about", about)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/cst", cst)

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	http.Handle("/", router)

	//cssHandler := http.FileServer(http.Dir("./css/"))
	//imagesHandler := http.FileServer(http.Dir("./images/"))

	//http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	//http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	port := ":3000"
	fmt.Println("Listening on localhost" + port)

	log.Fatal(http.ListenAndServe(port, router))
}
