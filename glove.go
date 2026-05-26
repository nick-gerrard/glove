package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Game struct {
	Title    string `json:"title"`
	Tagline  string `json:"tagline"`
	Version  string `json:"version"`
	GamePath string `json:"game_path"`
}

var games = []Game{}

//
// 	if err != nil {
// 		http.Error(w, "template parse error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
// 		log.Printf("template execute error: %v", err)
// 		return
// 	}
// }

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", games)
}

// func gameHandler(w http.ResponseWriter, r *http.Request) {
// 	renderPage(w, nil, "templates/game.html")
// }
//
// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	renderPage(w, nil, "templates/about.html")
// }

func main() {
	data, err := os.ReadFile("static/games.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &games)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/game", gameHandler)
	// http.HandleFunc("/about", aboutHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
