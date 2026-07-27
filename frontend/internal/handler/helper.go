package handler

import (
	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
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
