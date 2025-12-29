package main

import (
	ascii "asciiart/features"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tmpl *template.Template
var notFoundTmpl *template.Template
var internalErrorTmpl *template.Template
var badRequestTmpl *template.Template

type PageData struct {
	Input  string
	Banner string
	Art    string
}

func init() {
	var err error
	tmpl, err = template.ParseFiles("./template/index.html")
	if err != nil {
		log.Fatal("Could not pre-parse templates: ", err)
	}

	// Parse 404 template
	notFoundTmpl, err = template.ParseFiles("./template/404.html")
	if err != nil {
		log.Fatal("Could not pre-parse 404 template: ", err)
	}
	// Parse 500 template
	internalErrorTmpl, err = template.ParseFiles("./template/500.html")
	if err != nil {
		log.Fatal("Could not pre-parse 500 template: ", err)
	}
	// Parse 400 template
	badRequestTmpl, err = template.ParseFiles("./template/400.html")
	if err != nil {
		log.Fatal("Could not pre-parse 400 template: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Serve custom 404 page
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl.Execute(w, nil)
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
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

		bannerData := ascii.ReadBanner(banner)
		if bannerData == "" {
			w.WriteHeader(http.StatusInternalServerError)
			internalErrorTmpl.Execute(w, nil)
			return
		}

		splitInput := strings.Split(strings.ReplaceAll(inputStr, "\r\n", "\n"), "\n")
		sliceBanner := strings.Split(bannerData, "\n")

		art, err := ascii.DrawingInput(splitInput, sliceBanner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			internalErrorTmpl.Execute(w, nil)
			return
		}
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
