package main

import (
	ascii "asciiart/features"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type PageData struct {
	Input  string
	Banner string
	Art    string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	Input := r.FormValue("text")
	Banner := r.FormValue("banner")
	//OutPut := ""
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
		splitInput := strings.Split(input, "\\n")
		fmt.Println(splitInput)
		fmt.Println(banner)
		err = tmpl.Execute(w, PageData{Input: Input, Banner: Banner})
		bannerData := ascii.ReadBanner(banner)
		fmt.Println(bannerData)
	}
	// bannerSlice := strings.Split(bannerData, "\n")
	// art, err := ascii.DrawingInput(splitInput, bannerSlice)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// OutPut = art

	// err = tmpl.Execute(w, PageData{Input: Input, Banner: Banner, Art: OutPut})
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

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
