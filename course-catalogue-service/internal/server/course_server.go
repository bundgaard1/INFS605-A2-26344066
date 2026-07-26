package server

import (
	"context"
	"log"

	coursecatalogue "osbourne.local/course-catalogue-service/gen/course-catalogue"
	"osbourne.local/course-catalogue-service/internal/domain"
	"osbourne.local/course-catalogue-service/internal/service"
)

type CourseServer struct {
	coursecatalogue.UnimplementedCourseCatalogueServiceServer
	courseSvc *service.CourseService
}

func NewCourseServer(courseSvc *service.CourseService) *CourseServer {
	return &CourseServer{
		courseSvc: courseSvc,
	}
}

func (s *CourseServer) GetCourse(ctx context.Context, req *coursecatalogue.GetCourseRequest) (*coursecatalogue.GetCourseResponse, error) {
	log.Printf("Received gRPC request for course_id: %s", req.GetCourseId())

	course, err := s.courseSvc.GetCourse(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}

	courseProto := toProtoCourse(course)

	return &coursecatalogue.GetCourseResponse{
		Course: courseProto,
	}, nil
}

func toProtoCourse(course *domain.Course) *coursecatalogue.Course {
	return &coursecatalogue.Course{
		Id:          course.ID,
		Code:        course.Code,
		Title:       course.Title,
		Description: course.Description,
		Credits:     int32(course.Credits),
	}
}
