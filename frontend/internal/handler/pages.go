package handler

import (
	"log"
	"net/http"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	"osbourne.local/frontend/gen/notification"
	"osbourne.local/frontend/internal/domain"
	"osbourne.local/frontend/internal/view"
)

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	currentUser := UserFromContext(r.Context())

	res, err := h.clients.CourseCatalogue.Client.ListEnrolledCourses(
		r.Context(),
		&coursecatalogue.ListEnrolledCoursesRequest{
			UserId: currentUser.ID,
		},
	)
	if err != nil {
		http.Error(w, "Could not fetch the user's enrollments", http.StatusBadGateway)
		log.Printf("gRPC call ListEnrolledCourses failed: %v", err)
		return
	}

	courses := make([]domain.Course, 0, len(res.GetEnrolledCourses()))
	for _, c := range res.GetEnrolledCourses() {
		courses = append(courses, toDomainCourse(c))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.DashboardPage(currentUser, courses).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.ProfilePage(currentUser).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	res, err := h.clients.Notification.Client.GetUserNotifications(
		r.Context(),
		&notification.NotificationsRequest{UserId: user.ID},
	)
	if err != nil {
		http.Error(w, "Could not fetch notifications", http.StatusBadGateway)
		log.Printf("gRPC call GetUserNotifications failed: %v", err)
		return
	}

	notifications := make([]domain.Notification, 0, len(res.GetNotifications()))
	for _, n := range res.GetNotifications() {
		notifications = append(notifications, toDomainNotification(n))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.NotificationsPage(user, notifications).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCourseCatalog(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	res, err := h.clients.CourseCatalogue.Client.ListCourses(
		r.Context(),
		&coursecatalogue.ListCoursesRequest{
			Page:     1,
			PageSize: 10,
		},
	)
	if err != nil {
		http.Error(w, "Could not fetch course catalog", http.StatusBadGateway)
		log.Printf("gRPC call ListCourses failed: %v", err)
		return
	}

	courses := make([]domain.Course, 0, len(res.GetCourses()))
	for _, course := range res.GetCourses() {
		courses = append(courses, toDomainCourse(course))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.CatalogPage(user, courses).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCoursePage(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	courseID := r.PathValue("id")

	res, err := h.clients.CourseCatalogue.Client.GetCourse(
		r.Context(),
		&coursecatalogue.GetCourseRequest{CourseId: courseID},
	)
	if err != nil {
		http.Error(w, "Could not fetch course details", http.StatusBadGateway)
		log.Printf("gRPC call GetCourse failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.CoursePage(user, toDomainCourse(res.GetCourse())).Render(r.Context(), w); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
