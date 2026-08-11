package server

import (
	"context"

	coursecontent "osbourne.local/course-content-service/gen/course-content"
	"osbourne.local/course-content-service/internal/domain"
	"osbourne.local/course-content-service/internal/service"
)

type ContentServer struct {
	coursecontent.CourseContentServiceServer
	sourceSvc *service.ModuleService
}

func NewContentServer(sourceSvc *service.ModuleService) *ContentServer {
	return &ContentServer{
		sourceSvc: sourceSvc,
	}
}

func (s *ContentServer) Create(ctx context.Context, req *coursecontent.CreateModuleRequest) (*coursecontent.CreateModuleResponse, error) {
	module := &domain.Module{
		ID:       "",
		CourseID: req.GetCourseId(),
		Title:    req.GetTitle(),
		Text:     req.GetText(),
	}

	err := s.sourceSvc.CreateModule(ctx, module)
	if err != nil {
		return nil, err
	}

	m, err := s.sourceSvc.GetModule(ctx, module.ID)
	if err != nil {
		return nil, err
	}

	moduleProto := toProtoModule(m)

	return &coursecontent.CreateModuleResponse{
		Module: moduleProto,
	}, nil
}

func (s *ContentServer) GetModule(ctx context.Context, req *coursecontent.GetModuleRequest) (*coursecontent.GetModuleResponse, error) {
	module, err := s.sourceSvc.GetModule(ctx, req.GetModuleId())
	if err != nil {
		return nil, err
	}

	moduleProto := toProtoModule(module)

	return &coursecontent.GetModuleResponse{
		Module: moduleProto,
	}, nil
}

func (s *ContentServer) UpdateModule(ctx context.Context, req *coursecontent.UpdateModuleRequest) (*coursecontent.UpdateModuleResponse, error) {
	update := &service.UpdateModuleInput{
		ID:    req.GetId(),
		Title: req.GetTitle(),
		Text:  req.GetText(),
	}

	err := s.sourceSvc.UpdateModule(ctx, update)
	if err != nil {
		return nil, err
	}

	m, err := s.sourceSvc.GetModule(ctx, update.ID)
	if err != nil {
		return nil, err
	}

	moduleProto := toProtoModule(m)

	return &coursecontent.UpdateModuleResponse{
		Module: moduleProto,
	}, nil
}

func (s *ContentServer) DeleteModule(ctx context.Context, req *coursecontent.DeleteModuleRequest) (*coursecontent.DeleteModuleResponse, error) {
	err := s.sourceSvc.DeleteModule(ctx, req.GetModuleId())
	if err != nil {
		return nil, err
	}

	return &coursecontent.DeleteModuleResponse{}, nil
}

func (s *ContentServer) ListModulesByCourseID(ctx context.Context, req *coursecontent.ListModulesByCourseIDRequest) (*coursecontent.ListModulesByCourseIDResponse, error) {
	modules, err := s.sourceSvc.ListModulesByCourseID(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}

	var moduleProtos []*coursecontent.Module
	for _, module := range modules {
		moduleProtos = append(moduleProtos, toProtoModule(module))
	}

	return &coursecontent.ListModulesByCourseIDResponse{
		Modules: moduleProtos,
	}, nil
}

func toProtoModule(module *domain.Module) *coursecontent.Module {
	return &coursecontent.Module{
		Id:       module.ID,
		CourseId: module.CourseID,
		Title:    module.Title,
		Text:     module.Text,
	}
}
