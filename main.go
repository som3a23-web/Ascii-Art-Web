package main

import (
	ascii "asciiart/features"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tmpl *template.Template

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
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
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

	r.Body = http.MaxBytesReader(w, r.Body, 10240)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "413 Request Entity Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	inputStr := r.FormValue("text")
	banner := r.FormValue("banner")

	if inputStr == "" || banner == "" {
		http.Error(w, "400 Bad Request: Missing text or banner", http.StatusBadRequest)
		return
	}

	if !isValidBanner(banner) {
		http.Error(w, "400 Bad Request: Invalid banner", http.StatusBadRequest)
		return
	}

	if !isValidASCII(inputStr) {
		http.Error(w, "400 Bad Request: Input contains non-ASCII characters", http.StatusBadRequest)
		return
	}

	bannerData := ascii.ReadBanner(banner)
	if bannerData == "" {
		http.Error(w, "500 Internal Server Error: Banner file missing", http.StatusInternalServerError)
		return
	}

	splitInput := strings.Split(strings.ReplaceAll(inputStr, "\r\n", "\n"), "\n")
	sliceBanner := strings.Split(bannerData, "\n")

	art, err := ascii.DrawingInput(splitInput, sliceBanner)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Input:  inputStr,
		Banner: banner,
		Art:    art,
	}
	tmpl.Execute(w, data)
}

func isValidBanner(banner string) bool {
	return banner == "standard" || banner == "shadow" || banner == "thinkertoy"
}

// Helper to check for non-ASCII characters (range 32-126 + newlines)
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
