package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Input  string
	Banner string
	Art    string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	Input := r.FormValue("textInput")
	Banner := r.FormValue("banner")
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{Input: Input, Banner: Banner})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// func checkErr(err error) {
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
