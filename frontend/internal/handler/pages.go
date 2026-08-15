package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	coursecontent "osbourne.local/frontend/gen/course-content"
	"osbourne.local/frontend/gen/notification"
	"osbourne.local/frontend/internal/view"
)

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())

	res, err := h.clients.CourseCatalogue.Client.ListEnrolledCourses(
		r.Context(),
		&coursecatalogue.ListEnrolledCoursesRequest{
			UserId: currentUser.ID,
		},
	)
	if err != nil {
		fetchError(w, "Could not fetch the user's enrollments", err)
		return
	}

	renderPage(w, r, view.DashboardPage(view.DashboardPageData{
		PageData: view.PageData{User: currentUser},
		Courses:  toDomainCourses(res.GetEnrolledCourses()),
	}))

}

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())
	renderPage(w, r, view.ProfilePage(view.ProfilePageData{
		PageData: view.PageData{User: currentUser},
	}))
}

func (h *Handler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	res, err := h.clients.Notification.Client.GetUserNotifications(
		r.Context(),
		&notification.NotificationsRequest{UserId: user.ID},
	)
	if err != nil {
		fetchError(w, "Could not fetch notifications", err)
		return
	}

	renderPage(w, r, view.NotificationsPage(view.NotificationsPageData{
		PageData:      view.PageData{User: user},
		Notifications: toDomainNotifications(res.GetNotifications()),
	}))
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
		fetchError(w, "Could not fetch course catalog", err)
		return
	}

	renderPage(w, r, view.CatalogPage(view.CatalogPageData{
		PageData: view.PageData{User: user},
		Courses:  toDomainCourses(res.GetCourses()),
	}))
}

func (h *Handler) HandleCoursePage(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	courseID := chi.URLParam(r, "courseID")

	res, err := h.clients.CourseCatalogue.Client.GetCourse(
		r.Context(),
		&coursecatalogue.GetCourseRequest{CourseId: courseID},
	)

	if err != nil {
		fetchError(w, "Could not fetch course details", err)
		return
	}

	res2, err := h.clients.CourseContent.Client.ListModulesByCourseID(
		r.Context(),
		&coursecontent.ListModulesByCourseIDRequest{CourseId: courseID},
	)

	if err != nil {
		fetchError(w, "Could not fetch course modules", err)
		return
	}

	modules := toDomainModules(res2.GetModules())

	coursePageData := view.CoursePageData{
		PageData: view.PageData{User: user},
		Course:   toDomainCourse(res.GetCourse()),
		Modules:  modules,
	}

	renderPage(w, r, view.CoursePage(coursePageData))
}
