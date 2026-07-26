package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"osbourne.local/frontend/gen/notification"
	"osbourne.local/frontend/gen/profile"
	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/ui"
)

type userContextKey string

const userKey userContextKey = "currentUser"

type App struct {
	tmpls              map[string]*template.Template
	profileClient      *grpcclient.ProfileClient
	notificationClient *grpcclient.NotificationClient
}

// Global template data struktur – gør det typsikkert at sende data til 'base'
type PageData struct {
	User map[string]string
	Data interface{}
}

func main() {
	tmplDashboard, err := template.ParseFS(ui.Files, "templates/layouts/base.html", "templates/dashboard.html")
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates: %v", err)
	}
	tmplProfile, err := template.ParseFS(ui.Files, "templates/layouts/base.html", "templates/profile.html")
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates: %v", err)
	}
	tmplNotifications, err := template.ParseFS(ui.Files, "templates/layouts/base.html", "templates/notifications.html")
	if err != nil {
		log.Fatalf("Fejl ved parsing af templates: %v", err)
	}

	profileServiceAddr := "dns:///profile-service:50051"
	profileClient, err := grpcclient.NewProfileClient(profileServiceAddr)
	if err != nil {
		log.Fatalf("Fejl ved oprettelse af gRPC-klient: %v", err)
	}
	defer profileClient.Close()

	notificationServiceAddr := "dns:///notification-service:50052"
	notificationClient, err := grpcclient.NewNotificationClient(notificationServiceAddr)
	if err != nil {
		log.Fatalf("Fejl ved oprettelse af gRPC-klient: %v", err)
	}
	defer notificationClient.Close()

	app := &App{
		tmpls: map[string]*template.Template{
			"dashboard":     tmplDashboard,
			"profile":       tmplProfile,
			"notifications": tmplNotifications,
		},
		profileClient:      profileClient,
		notificationClient: notificationClient,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(ui.Files)))

	// Wrap dine handlers i authenticatedUser middleware
	mux.Handle("/", app.authenticatedUser(http.HandlerFunc(app.handleDashboard)))
	mux.Handle("/profile", app.authenticatedUser(http.HandlerFunc(app.handleProfile)))
	mux.Handle("/notifications", app.authenticatedUser(http.HandlerFunc(app.handleNotifications)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Go HTML Template Web Service running on port :" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func (app *App) authenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("id")
		if userID == "" {
			userID = "12345"
		}

		res, err := app.profileClient.Client.GetUserProfile(
			r.Context(),
			&profile.ProfileRequest{UserId: userID},
		)

		var userMap map[string]string
		if err == nil {
			userMap = map[string]string{
				"Id":   res.GetId(),
				"Name": res.GetName(),
				"Role": res.GetRole(),
			}
		} else {
			// Fallback i tilfælde af gRPC fejl, så UI ikke crasher
			log.Printf("gRPC Hentning af nav-bruger fejlede: %v", err)
			userMap = map[string]string{
				"Name": "Gæst",
				"Role": "Ukendt",
			}
		}

		ctx := context.WithValue(r.Context(), userKey, userMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper til at trække brugeren ud af context
func getUserFromContext(ctx context.Context) map[string]string {
	if u, ok := ctx.Value(userKey).(map[string]string); ok {
		return u
	}
	return map[string]string{"Name": "Gæst", "Role": "Ukendt"}
}

func (app *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Pak data ind i PageData, så .User automatisk er tilgængelig for base.html
	data := PageData{
		User: getUserFromContext(r.Context()),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.tmpls["dashboard"].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) handleProfile(w http.ResponseWriter, r *http.Request) {
	// Hent den opslåede bruger
	userID := r.URL.Query().Get("id")
	if userID == "" {
		userID = "12345"
	}

	res, err := app.profileClient.Client.GetUserProfile(
		r.Context(),
		&profile.ProfileRequest{UserId: userID},
	)

	if err != nil {
		http.Error(w, "Kunne ikke hente profil", http.StatusBadGateway)
		log.Printf("gRPC-kald GetUserProfile fejlede: %v", err)
		return
	}

	profileUser := map[string]string{
		"Id":   res.GetId(),
		"Name": res.GetName(),
		"Role": res.GetRole(),
	}

	// Send BÅDE navbarens bruger (.User) og sidens specifikke data (.Data)
	pageData := PageData{
		User: getUserFromContext(r.Context()),
		Data: profileUser,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.tmpls["profile"].ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) handleNotifications(w http.ResponseWriter, r *http.Request) {
	// Hent den opslåede bruger
	userID := r.URL.Query().Get("id")
	if userID == "" {
		userID = "12345"
	}

	res, err := app.notificationClient.Client.GetUserNotifications(
		r.Context(),
		&notification.NotificationsRequest{UserId: userID},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente notifikationer", http.StatusBadGateway)
		log.Printf("gRPC-kald GetUserNotifications fejlede: %v", err)
		return
	}

	notifications := res.GetNotifications()

	PageData := PageData{
		User: getUserFromContext(r.Context()),
		Data: map[string]interface{}{
			"Notifications": notifications,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.tmpls["notifications"].ExecuteTemplate(w, "base", PageData); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
