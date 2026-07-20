package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	http.HandleFunc("/", handleDashboard)

	log.Println("Go HTML Template Web Service running on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"User": map[string]string{
			"Name": "Andy",
			"Role": "student",
		},
		"Courses": []map[string]string{
			{"ID": "INFS-605", "Title": "Microservices Programming Project"},
			{"ID": "COMP-901", "Title": "Distributed Systems"},
		},
	}

	lp := filepath.Join("templates", "layouts", "base.html")
	fp := filepath.Join("templates", "dashboard.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
