package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Game struct {
	Title    string `json:"title"`
	Tagline  string `json:"tagline"`
	Version  string `json:"version"`
	GamePath string `json:"game_path"`
}

var games = []Game{}

func addHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", games)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game := r.PathValue("title")
	for _, g := range games {
		if strings.EqualFold(game, g.Title) {
			tmpl, err := template.ParseFiles("templates/base.html", "templates/game.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "base", g)
		}
	}
}

func main() {
	data, err := os.ReadFile("static/games.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &games)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/games/{title}", gameHandler)
	log.Fatal(http.ListenAndServe(":8000", addHeaders(mux)))
}
