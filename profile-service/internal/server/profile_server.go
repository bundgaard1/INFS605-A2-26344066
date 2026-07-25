package server

import (
	"context"
	"log"

	"osbourne.local/profile-service/gen/profile"
	"osbourne.local/profile-service/internal/service"
)

type ProfileServer struct {
	profile.UnimplementedProfileServiceServer
	profileSvc *service.ProfileService
}

func NewProfileServer(profileSvc *service.ProfileService) *ProfileServer {
	return &ProfileServer{
		profileSvc: profileSvc,
	}
}

func (s *ProfileServer) GetUserProfile(ctx context.Context, req *profile.ProfileRequest) (*profile.ProfileResponse, error) {
	log.Printf("Modtog gRPC-forespørgsel for profile_id: %s", req.GetUserId())

	return s.profileSvc.GetProfile(ctx, req.GetUserId())
}
