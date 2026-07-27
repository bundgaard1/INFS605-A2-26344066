package handler

import (
	"log"
	"net/http"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	"osbourne.local/frontend/gen/notification"
	"osbourne.local/frontend/internal/domain"
)

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	currentUser := UserFromContext(r.Context())

	res, err := h.courseCatalogueClient.Client.ListEnrolledCourses(
		r.Context(),
		&coursecatalogue.ListEnrolledCoursesRequest{
			UserId: currentUser.ID,
		},
	)

	if err != nil {
		http.Error(w, "Kunne ikke hente brugerens tilmeldinger", http.StatusBadGateway)
		log.Printf("gRPC-kald ListEnrolledCourses fejlede: %v", err)
		return
	}

	enrolledCourses := make([]domain.Course, 0, len(res.GetEnrolledCourses()))
	for _, c := range res.GetEnrolledCourses() {
		enrolledCourses = append(enrolledCourses, toDomainCourse(c))
	}

	data := domain.PageData{
		User: currentUser,
		Data: map[string]any{
			"EnrolledCourses": enrolledCourses,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "dashboard", data); err != nil {
		http.Error(w, "Render fejl: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())

	pageData := domain.PageData{
		User: currentUser,
		Data: currentUser,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "profile", pageData); err != nil {
		http.Error(w, "Render fejl: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	res, err := h.notificationClient.Client.GetUserNotifications(
		r.Context(),
		&notification.NotificationsRequest{UserId: user.ID},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente notifikationer", http.StatusBadGateway)
		log.Printf("gRPC-kald GetUserNotifications fejlede: %v", err)
		return
	}

	domainNotifications := make([]domain.Notification, 0, len(res.GetNotifications()))
	for _, n := range res.GetNotifications() {
		domainNotifications = append(domainNotifications, toDomainNotification(n))
	}

	pageData := domain.PageData{
		User: user,
		Data: map[string]any{
			"Notifications": domainNotifications,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "notifications", pageData); err != nil {
		http.Error(w, "Render fejl: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCourseCatalog(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	res, err := h.courseCatalogueClient.Client.ListCourses(
		r.Context(),
		&coursecatalogue.ListCoursesRequest{},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente kursus katalog", http.StatusBadGateway)
		log.Printf("gRPC-kald ListCourses fejlede: %v", err)
		return
	}

	domainCourses := make([]domain.Course, 0, len(res.GetCourses()))
	for _, course := range res.GetCourses() {
		domainCourses = append(domainCourses, toDomainCourse(course))
	}

	pageData := domain.PageData{
		User: user,
		Data: map[string]any{
			"Courses": domainCourses,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "course-catalog", pageData); err != nil {
		http.Error(w, "Render fejl: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCoursePage(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	courseID := r.PathValue("id")

	res, err := h.courseCatalogueClient.Client.GetCourse(
		r.Context(),
		&coursecatalogue.GetCourseRequest{CourseId: courseID},
	)
	if err != nil {
		http.Error(w, "Kunne ikke hente kursusdetaljer", http.StatusBadGateway)
		log.Printf("gRPC-kald GetCourse fejlede: %v", err)
		return
	}

	course := toDomainCourse(res.GetCourse())

	pageData := domain.PageData{
		User: user,
		Data: map[string]any{
			"Course": course,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "course-page", pageData); err != nil {
		http.Error(w, "Render fejl: "+err.Error(), http.StatusInternalServerError)
	}
}
