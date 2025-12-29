// package main

// import (
// 	ascii "asciiart/features"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"strings"
// )

// type PageData struct {
// 	Input  string
// 	Banner string
// 	Art    string
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	// Path Guard
// 	if r.URL.Path != "/" {
// 		http.Error(w, "404 Not Found", http.StatusNotFound)
// 		return
// 	}
// 	// Methods Guard
// 	if r.Method != http.MethodGet && r.Method != http.MethodPost {
// 		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	// Limit Rq body size and Parsing
// 	r.Body = http.MaxBytesReader(w, r.Body, 10240)
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "413 Request body too large (max 10KB)", http.StatusRequestEntityTooLarge)
// 		return
// 	}

// 	inputStr := r.FormValue("text")
// 	banner := r.FormValue("banner")
// 	outputArt := ""

// 	if r.Method == http.MethodPost {
// 		// Validation
// 		if inputStr == "" || banner == "" {
// 			http.Error(w, "400 Bad Request: Missing data", http.StatusBadRequest)
// 			return
// 		}

// 		if !isValidBanner(banner) {
// 			http.Error(w, "400 Bad Request: Unknown banner", http.StatusBadRequest)
// 			return
// 		}
// 		bannerData := ascii.ReadBanner(banner)
// 		// Handle missing file
// 		if bannerData == "" {
// 			http.Error(w, "500 Internal Server Error: Could not read banner file", http.StatusInternalServerError)
// 			return
// 		}
// 		splitInput := strings.Split(inputStr, "\n")
// 		sliceBanner := strings.Split(bannerData, "\n")
// 		art, err := ascii.DrawingInput(splitInput, sliceBanner)
// 		if err != nil {
// 			http.Error(w, "500 Internal Server Error: Error generating ASCII art", http.StatusInternalServerError)
// 			return
// 		}
// 		outputArt = art
// 	}

// 	// Render
// 	tmpl, err := template.ParseFiles("./template/index.html")
// 	if err != nil {
// 		http.Error(w, "500 Internal Server Error: Template missing", http.StatusInternalServerError)
// 		return
// 	}
// 	err = tmpl.Execute(w, PageData{Input: inputStr, Banner: banner, Art: outputArt})
// }

// func isValidBanner(banner string) bool {
// 	return banner == "standard" || banner == "shadow" || banner == "thinkertoy"
// }

// func main() {
// 	fs := http.FileServer(http.Dir("static"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))
// 	http.HandleFunc("/", homeHandler)
// 	log.Println("Server starting on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))

// }

package main

import (
	ascii "asciiart/features"

	"html/template"

	"log"

	"net/http"

	"strings"
)

type PageData struct {
	Input string

	Banner string

	Art string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	Input := r.FormValue("text")

	Banner := r.FormValue("banner")

	OutPut := ""

	if r.URL.Path != "/" {

		http.NotFound(w, r)

		return

	}

	tmpl, err := template.ParseFiles("./template/index.html")

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return

	}

	if Input != "" && Banner != "" {

		Args := []string{Input, Banner}

		input, banner := ascii.StoreInputAndBanner(Args)

		splitInput := strings.Split(input, "\n")

		bannerData := ascii.ReadBanner(banner)

		sliceBanner := strings.Split(bannerData, "\n")

		art, err := ascii.DrawingInput(splitInput, sliceBanner)

		OutPut = art

		checkErr(err)
	}

	err = tmpl.Execute(w, PageData{Input: Input, Banner: Banner, Art: OutPut})

}

func main() {

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	log.Println("Server starting on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func checkErr(err error) {

	if err != nil {
		log.Fatal(err)
	}

}
