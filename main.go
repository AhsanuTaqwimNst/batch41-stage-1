package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	route.HandleFunc("/send-data-add-project", sendDataAddProject).Methods("POST")

	fmt.Println("Server running at localhost port 8000")

	http.ListenAndServe("localhost:8000", route)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func sendDataAddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("project-name")
	inputStart := r.PostForm.Get("input-start")
	inputDate := r.PostForm.Get("input-date")
	description := r.PostForm.Get("description")
	var techno []string
	techno = r.Form["techno"]
	uploadImg := r.PostForm.Get("upload-img")

	fmt.Println("Project Name : " + projectName)
	fmt.Println("Input Start : " + inputStart)
	fmt.Println("End Date : " + inputDate)
	fmt.Println("Description : " + description)
	fmt.Println(techno)
	fmt.Println("Upload Image : " + uploadImg)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
