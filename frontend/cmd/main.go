package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"osbourne.local/frontend/gen/profile"
	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/ui"
)

type App struct {
	tmpl          *template.Template
	profileClient *grpcclient.ProfileClient
}

func main() {

	// Setup
	tmpl, err := template.ParseFS(ui.Files, "templates/layouts/base.html", "templates/dashboard.html")
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates: %v", err)
	}

	profileServiceAddr := "dns:///profile-service:50051"
	profileClient, err := grpcclient.NewProfileClient(profileServiceAddr)
	if err != nil {
		log.Fatalf("Fejl ved oprettelse af gRPC-klient: %v", err)
	}
	defer profileClient.Close()

	app := &App{
		tmpl:          tmpl,
		profileClient: profileClient,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(ui.Files)))

	mux.HandleFunc("/", app.handleDashboard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Go HTML Template Web Service running on port :" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func (app *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := app.profileClient.Client.GetUserProfile(
		r.Context(),
		&profile.ProfileRequest{
			UserId: "12345",
		},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente profil", http.StatusBadGateway)
		log.Printf("gRPC-kald GetUserProfile fejlede: %v", err)
		return
	}

	data := map[string]interface{}{
		"User": map[string]string{
			"Name": res.GetName(),
			"Role": res.GetRole(),
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := app.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
