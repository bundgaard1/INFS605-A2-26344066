package handler

import (
	"encoding/json"
	"net/http"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
)

func (h *Handler) HandleEnrollCourse(w http.ResponseWriter, r *http.Request) {
	// 1. Read HTTP request from JS/HTMX
	var req struct {
		CourseID string `json:"course_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Get the authenticated user (e.g. from session/JWT cookie)
	userID := UserFromContext(r.Context()).ID

	// 3. Call the gRPC service internally in Go
	resp, err := h.courseCatalogueClient.Client.EnrollUser(r.Context(),
		&coursecatalogue.EnrollUserRequest{
			UserId:   userID,
			CourseId: req.CourseID,
		})

	if err != nil {
		http.Error(w, "Could not complete enrollment", http.StatusInternalServerError)
		return
	}

	// 4. Send clean HTTP/JSON back to the browser
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": resp.GetSuccess(),
	})
}
