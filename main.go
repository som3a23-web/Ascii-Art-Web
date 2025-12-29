package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title  string
	Body   string
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
	fmt.Println(Input)
	fmt.Print(Banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	// stringToArt, banner := ascii.StoreInputAndBanner()
	// bannerFile := ascii.ReadBanner(banner)
	// bannerSlice := strings.Split(bannerFile, "\n")
	// splitInput := strings.Split(stringToArt, "\\n")
	// outPut := ascii.DrawingInput(splitInput, bannerSlice)

}

// func checkErr(err error) {
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
