package handler

import (
	"encoding/json"
	"net/http"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
)

func (h *Handler) HandleEnrollCourse(w http.ResponseWriter, r *http.Request) {
	// 1. Læs HTTP request fra JS/HTMX
	var req struct {
		CourseID string `json:"course_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ugyldig request body", http.StatusBadRequest)
		return
	}

	// 2. Hent den autentificerede bruger (f.eks. fra session/JWT cookie)
	userID := UserFromContext(r.Context()).ID

	// 3. Kalder gRPC servicen internt i Go
	resp, err := h.courseCatalogueClient.Client.EnrollUser(r.Context(),
		&coursecatalogue.EnrollUserRequest{
			UserId:   userID,
			CourseId: req.CourseID,
		})

	if err != nil {
		http.Error(w, "Kunne ikke gennemføre tilmelding", http.StatusInternalServerError)
		return
	}

	// 4. Send ren HTTP/JSON tilbage til browseren
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": resp.GetSuccess(),
	})
}
