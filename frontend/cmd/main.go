package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	usersServiceAddr := os.Getenv("USERS_SERVICE_ADDR")
	if usersServiceAddr == "" {
		usersServiceAddr = "localhost:50051"
	}

	conn, err := grpc.NewClient(usersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Kunne ikke oprette forbindelse til users-service", http.StatusBadGateway)
		log.Printf("Kunne ikke oprette gRPC-forbindelse til %s: %v", usersServiceAddr, err)
		return
	}
	defer conn.Close()

	client := student.NewStudentServiceClient(conn)

	// Kalder gRPC servicen
	res, err := client.GetStudentProfile(context.Background(), &student.StudentRequest{
		StudentId: "12345",
	})
	if err != nil {
		http.Error(w, "Kunne ikke hente studentprofil", http.StatusBadGateway)
		log.Printf("gRPC-kald GetStudentProfile fejlede mod %s: %v", usersServiceAddr, err)
		return
	}
	log.Println("Navn modtaget fra gRPC:", res.GetName())

	data := map[string]interface{}{
		"User": map[string]string{
			"Name": res.GetName(),
			"Role": res.GetRole(),
		},
		"Courses": res.GetCourses(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Eksekver den allerede parsede skabelon direkte fra hukommelsen
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
