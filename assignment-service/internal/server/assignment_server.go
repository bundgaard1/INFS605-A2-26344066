package server

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	assignmentpb "osbourne.local/assignment-service/gen/assignment"
	"osbourne.local/assignment-service/internal/domain"
	"osbourne.local/assignment-service/internal/service"
)

type AssignmentServer struct {
	assignmentpb.AssignmentServiceServer
	svc *service.AssignmentService
}

func NewAssignmentServer(svc *service.AssignmentService) *AssignmentServer {
	return &AssignmentServer{
		svc: svc,
	}
}

func (s *AssignmentServer) CreateAssignment(ctx context.Context, req *assignmentpb.CreateAssignmentRequest) (*assignmentpb.CreateAssignmentResponse, error) {
	assignment := &domain.Assignment{
		CourseID:    req.GetCourseId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}
	if req.GetDueDate() != nil {
		assignment.DueDate = req.GetDueDate().AsTime()
	}

	if err := s.svc.CreateAssignment(ctx, assignment); err != nil {
		return nil, err
	}

	return &assignmentpb.CreateAssignmentResponse{
		Assignment: toProtoAssignment(assignment),
	}, nil
}

func (s *AssignmentServer) GetAssignment(ctx context.Context, req *assignmentpb.GetAssignmentRequest) (*assignmentpb.GetAssignmentResponse, error) {
	assignment, err := s.svc.GetAssignment(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, status.Error(codes.NotFound, "assignment not found")
	}

	return &assignmentpb.GetAssignmentResponse{
		Assignment: toProtoAssignment(assignment),
	}, nil
}

func (s *AssignmentServer) UploadSubmission(stream assignmentpb.AssignmentService_UploadSubmissionServer) error {
	first, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "expected metadata as first message: %v", err)
	}

	meta := first.GetMetadata()
	if meta == nil {
		return status.Error(codes.InvalidArgument, "first message must contain metadata")
	}
	if meta.GetAssignmentId() == "" || meta.GetStudentId() == "" || meta.GetFilename() == "" {
		return status.Error(codes.InvalidArgument, "metadata must contain assignment_id, student_id and filename")
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		for {
			req, rerr := stream.Recv()
			if rerr == io.EOF {
				return
			}
			if rerr != nil {
				pw.CloseWithError(rerr)
				return
			}
			if _, werr := pw.Write(req.GetChunk()); werr != nil {
				pw.CloseWithError(werr)
				return
			}
		}
	}()

	submission, err := s.svc.SubmitAssignment(stream.Context(), service.SubmitAssignmentInput{
		AssignmentID: meta.GetAssignmentId(),
		StudentID:    meta.GetStudentId(),
		FileName:     meta.GetFilename(),
		Size:         meta.GetSize(),
	}, pr)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&assignmentpb.UploadSubmissionResponse{
		SubmissionId: submission.ID,
		FileUrl:      submission.FileURL,
		FileName:     submission.FileName,
		FileSize:     submission.FileSize,
	})
}

func toProtoAssignment(assignment *domain.Assignment) *assignmentpb.Assignment {
	return &assignmentpb.Assignment{
		Id:          assignment.ID,
		CourseId:    assignment.CourseID,
		Title:       assignment.Title,
		Description: assignment.Description,
		DueDate:     timestamppb.New(assignment.DueDate),
	}
}
