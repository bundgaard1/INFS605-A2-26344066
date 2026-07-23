package main

import (
	"html/template"
	"log"
	"net/http"

	"osbourne.local/frontend/gen/student"
	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/ui"
)

type App struct {
	tmpl          *template.Template
	studentClient *grpcclient.StudentClient
}

func main() {

	// Setup
	tmpl, err := template.ParseFS(ui.Files, "templates/layouts/base.html", "templates/dashboard.html")
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates: %v", err)
	}

	usersServiceAddr := "dns:///profile-service:50051"
	studentClient, err := grpcclient.NewStudentClient(usersServiceAddr)
	if err != nil {
		log.Fatalf("Fejl ved oprettelse af gRPC-klient: %v", err)
	}
	defer studentClient.Close()

	app := &App{
		tmpl:          tmpl,
		studentClient: studentClient,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(ui.Files)))

	mux.HandleFunc("/", app.handleDashboard)

	log.Println("Go HTML Template Web Service running on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func (app *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := app.studentClient.Client.GetStudentProfile(
		r.Context(),
		&student.StudentRequest{
			StudentId: "12345",
		},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente studentprofil", http.StatusBadGateway)
		log.Printf("gRPC-kald GetStudentProfile fejlede: %v", err)
		return
	}

	data := map[string]interface{}{
		"User": map[string]string{
			"Name": res.GetName(),
			"Role": res.GetRole(),
		},
		"Courses": res.GetCourses(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := app.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
