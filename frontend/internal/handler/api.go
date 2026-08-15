package handler

import (
	"encoding/json"
	"log"
	"net/http"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
)

type enrollResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *Handler) HandleEnrollCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.FormValue("course_id")
	if courseID == "" {
		writeJSON(w, http.StatusBadRequest, enrollResponse{Success: false, Message: "Missing course_id"})
		return
	}

	userID := UserFromContext(r.Context()).ID

	courseRes, courseErr := h.clients.CourseCatalogue.Client.GetCourse(
		r.Context(),
		&coursecatalogue.GetCourseRequest{CourseId: courseID},
	)

	_, enrollErr := h.clients.CourseCatalogue.Client.EnrollUser(r.Context(),
		&coursecatalogue.EnrollUserRequest{
			UserId:   userID,
			CourseId: courseID,
		})

	if enrollErr != nil {
		log.Printf("gRPC call EnrollUser failed: %v", enrollErr)
		writeJSON(w, grpcToHTTPStatus(enrollErr), enrollResponse{Success: false, Message: "Could not complete enrollment"})
		return
	}

	title := "Enrolled successfully"
	if courseErr == nil && courseRes.GetCourse() != nil {
		title = "Enrolled in " + courseRes.GetCourse().GetTitle()
	}

	writeJSON(w, http.StatusOK, enrollResponse{Success: true, Message: title})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}
