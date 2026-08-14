package handler

import (
	"context"
	"io/fs"
	"log"
	"net/http"

	"osbourne.local/frontend/gen/profile"
	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/internal/domain"
)

type contextKey string

const userKey contextKey = "currentUser"

type Handler struct {
	clients *grpcclient.Clients
}

func New(clients *grpcclient.Clients) *Handler {
	return &Handler{
		clients: clients,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, staticFiles fs.FS) {
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	mux.Handle("/", h.Authenticate(http.HandlerFunc(h.HandleDashboard)))
	mux.Handle("/profile", h.Authenticate(http.HandlerFunc(h.HandleProfile)))
	mux.Handle("/notifications", h.Authenticate(http.HandlerFunc(h.HandleNotifications)))
	mux.Handle("/course-catalog", h.Authenticate(http.HandlerFunc(h.HandleCourseCatalog)))
	mux.Handle("/courses/{id}", h.Authenticate(http.HandlerFunc(h.HandleCoursePage)))
	// API Endpoints
	mux.Handle("/api/courses/enroll", h.Authenticate(http.HandlerFunc(h.HandleEnrollCourse)))
}

// Middleware: Fetches the user once into the context
func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("id")
		if userID == "" {
			userID = "12345"
		}

		res, err := h.clients.Profile.Client.GetUserProfile(
			r.Context(),
			&profile.ProfileRequest{UserId: userID},
		)

		user := domain.User{Name: "Guest", Role: "Unknown"}
		if err == nil {
			user = domain.User{
				ID:   res.GetId(),
				Name: res.GetName(),
				Role: res.GetRole(),
			}
		} else {
			log.Printf("gRPC user fetch failed: %v", err)
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) domain.User {
	if u, ok := ctx.Value(userKey).(domain.User); ok {
		return u
	}
	return domain.User{Name: "Guest", Role: "Unknown"}
}
