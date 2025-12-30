package main

import (
	ascii "asciiart/features"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var tmpl *template.Template
var notFoundTmpl *template.Template
var internalErrorTmpl *template.Template
var badRequestTmpl *template.Template
var methodNotAllowedTmpl *template.Template

type PageData struct {
	Input  string
	Banner string
	Art    string
}

func init() {
	var err error
	tmpl, err = template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal("Could not pre-parse templates: ", err)
	}

	// Parse 404 template
	notFoundTmpl, err = template.ParseFiles("./templates/404.html")
	if err != nil {
		log.Fatal("Could not pre-parse 404 template: ", err)
	}
	// Parse 500 template
	internalErrorTmpl, err = template.ParseFiles("./templates/500.html")
	if err != nil {
		log.Fatal("Could not pre-parse 500 template: ", err)
	}
	// Parse 400 template
	badRequestTmpl, err = template.ParseFiles("./templates/400.html")
	if err != nil {
		log.Fatal("Could not pre-parse 400 template: ", err)
	}
	// Parse 405 template
	methodNotAllowedTmpl, err = template.ParseFiles("./templates/405.html")
	if err != nil {
		log.Fatal("Could not pre-parse 405 template: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		methodNotAllowedTmpl.Execute(w, nil)
		return
	}

	tmpl.Execute(w, nil)
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		methodNotAllowedTmpl.Execute(w, nil)
		return
	}

	inputStr := r.FormValue("text")
	banner := r.FormValue("banner")
	var output string

	if inputStr != "" || banner != "" {

		if !isValidBanner(banner) {
			w.WriteHeader(http.StatusBadRequest)
			badRequestTmpl.Execute(w, nil)
			return
		}

		if !isValidASCII(inputStr) {
			w.WriteHeader(http.StatusBadRequest)
			badRequestTmpl.Execute(w, nil)
			return
		}

		bannerSelected, err := os.ReadFile("banner/" + banner + ".txt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			internalErrorTmpl.Execute(w, nil)
			return
		}
		bannerSelectedConv := string(bannerSelected)
		if bannerSelectedConv == "" {
			w.WriteHeader(http.StatusInternalServerError)
			internalErrorTmpl.Execute(w, nil)
			return
		}
		replaceInput := strings.ReplaceAll(inputStr, "\r\n", "\n")
		splitInput := strings.Split(replaceInput, "\n")
		sliceBanner := strings.Split(bannerSelectedConv, "\n")

		art := ascii.DrawingInput(splitInput, sliceBanner)
		output = art
	}

	data := PageData{
		Input:  inputStr,
		Banner: banner,
		Art:    output,
	}
	tmpl.Execute(w, data)
}

func isValidBanner(banner string) bool {
	return banner == "standard" || banner == "shadow" || banner == "thinkertoy"
}

func isValidASCII(s string) bool {
	for _, r := range s {
		if (r < 32 || r > 126) && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server starting on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
