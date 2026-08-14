package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	"osbourne.local/frontend/internal/view"
)

func (h *Handler) HandleEnrollCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.FormValue("course_id")
	if courseID == "" {
		http.Error(w, "Missing course_id", http.StatusBadRequest)
		return
	}

	userID := UserFromContext(r.Context()).ID

	fmt.Println("Trying to enroll user ", userID, " in course ", courseID)

	courseRes, courseErr := h.clients.CourseCatalogue.Client.GetCourse(
		r.Context(),
		&coursecatalogue.GetCourseRequest{CourseId: courseID},
	)
	if courseErr != nil {
		log.Printf("gRPC call GetCourse failed: %v", courseErr)
	}

	_, enrollErr := h.clients.CourseCatalogue.Client.EnrollUser(r.Context(),
		&coursecatalogue.EnrollUserRequest{
			UserId:   userID,
			CourseId: courseID,
		})

	components := make([]templ.Component, 0, 2)
	if courseErr == nil {
		components = append(components, view.CourseCard(toDomainCourse(courseRes.GetCourse())))
	}

	if enrollErr != nil {
		log.Printf("gRPC call EnrollUser failed: %v", enrollErr)
		components = append(components, view.Toast("Could not complete enrollment", "error"))
	} else if courseErr == nil {
		components = append(components, view.Toast("Enrolled in "+courseRes.GetCourse().GetTitle(), "success"))
	} else {
		components = append(components, view.Toast("Enrolled successfully", "success"))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templ.Join(components...).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
