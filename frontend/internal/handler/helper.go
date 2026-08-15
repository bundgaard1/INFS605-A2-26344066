package handler

import (
	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	coursecontent "osbourne.local/frontend/gen/course-content"
	"osbourne.local/frontend/gen/notification"
	"osbourne.local/frontend/internal/domain"
)

func toDomainCourse(c *coursecatalogue.Course) domain.Course {
	if c == nil {
		return domain.Course{}
	}

	return domain.Course{
		ID:          c.GetId(),
		Code:        c.GetCode(),
		Title:       c.GetTitle(),
		Description: c.GetDescription(),
		Credits:     int(c.GetCredits()),
	}
}

func toDomainNotification(n *notification.Notification) domain.Notification {
	if n == nil {
		return domain.Notification{}
	}
	return domain.Notification{
		ID:        n.GetId(),
		Title:     n.GetTitle(),
		Message:   n.GetMsg(),
		Read:      n.GetIsRead(),
		Link:      n.GetLink(),
		Timestamp: n.GetTimestamp().AsTime(),
	}
}

func toDomainCourses(courses []*coursecatalogue.Course) []domain.Course {
	result := make([]domain.Course, 0, len(courses))
	for _, c := range courses {
		result = append(result, toDomainCourse(c))
	}
	return result
}

func toDomainNotifications(notifications []*notification.Notification) []domain.Notification {
	result := make([]domain.Notification, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toDomainNotification(n))
	}
	return result
}

func toDomainModules(modules []*coursecontent.Module) []domain.Module {
	result := make([]domain.Module, 0, len(modules))
	for _, m := range modules {
		result = append(result, domain.Module{
			ID:       m.GetId(),
			CourseID: m.GetCourseId(),
			Title:    m.GetTitle(),
			Text:     m.GetText(),
		})
	}
	return result
}
