package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"osbourne.local/frontend/gen/student"
)

// Global variabel til at holde de færdigparsede skabeloner i hukommelsen
var tmpl *template.Template

func init() {
	var err error

	// Læs og parse skabelonerne fra disken ÉN gang, når programmet starter
	lp := filepath.Join("templates", "layouts", "base.html")
	fp := filepath.Join("templates", "dashboard.html")

	tmpl, err = template.ParseFiles(lp, fp)
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates ved opstart: %v", err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

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

	data := &student.StudentResponse{
		Id:   "123",
		Name: "Andy",
		Role: "student",
		Courses: []*student.Course{
			{Id: "INFS-605", Title: "Microservices Programming Project"},
			{Id: "COMP-901", Title: "Distributed Systems"},
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Eksekver den allerede parsede skabelon direkte fra hukommelsen
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
