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

type projectData struct {
	Name    string
	Surname string
}

var pdata projectData

func testHandler(w http.ResponseWriter, r *http.Request) {

	data := pageData{Title: "Deon", Navbar: "hello"}
	tmpl := template.Must(template.ParseFiles("test.html"))
	err := tmpl.Execute(w, data)
	//err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		tmpl := template.Must(template.ParseFiles("project.html", "headernav.html"))
		err := tmpl.Execute(w, pdata)
		//err := tmpl.Execute(os.Stdout, data)
		if err != nil {
			panic(err)
		}
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		pdata.Name = r.FormValue("name")
		pdata.Surname = r.FormValue("surname")
		fmt.Println("Name = ", pdata.Name)
		fmt.Println("Surname = ", pdata.Surname)
		fmt.Println("data posted")
		http.Redirect(w, r, "/project", http.StatusSeeOther)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	fmt.Println("project executed")

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

	data.Uname = pdata.Name
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
	pdata.Name = "Deon"
	pdata.Surname = "Kloppers"

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/about", about)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/project", projectHandler)
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
