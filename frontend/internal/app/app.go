package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/internal/handler"
	"osbourne.local/frontend/internal/render"
	"osbourne.local/frontend/ui"
)

type Config struct {
	Port                       string
	ProfileServiceAddr         string
	NotificationServiceAddr    string
	CourseCatalogueServiceAddr string
	CourseContentServiceAddr   string
}

type App struct {
	cfg                   Config
	profileClient         *grpcclient.ProfileClient
	notificationClient    *grpcclient.NotificationClient
	courseCatalogueClient *grpcclient.CourseCatalogueClient
	courseContentClient   *grpcclient.CourseContentClient
	server                *http.Server
}

func NewApp(cfg Config) (*App, error) {
	// 1. gRPC Clients
	pClient, err := grpcclient.NewProfileClient(cfg.ProfileServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile client: %w", err)
	}

	nClient, err := grpcclient.NewNotificationClient(cfg.NotificationServiceAddr)
	if err != nil {
		pClient.Close()
		return nil, fmt.Errorf("failed to create notification client: %w", err)
	}

	ccClient, err := grpcclient.NewCourseCatalogueClient(cfg.CourseCatalogueServiceAddr)
	if err != nil {
		pClient.Close()
		nClient.Close()
		return nil, fmt.Errorf("failed to create course catalogue client: %w", err)
	}

	cContentClient, err := grpcclient.NewCourseContentClient(cfg.CourseContentServiceAddr)
	if err != nil {
		pClient.Close()
		nClient.Close()
		ccClient.Close()
		return nil, fmt.Errorf("failed to create course content client: %w", err)
	}

	// 2. Render & Handlers
	renderer, err := render.New(ui.Files)
	if err != nil {
		pClient.Close()
		nClient.Close()
		ccClient.Close()
		cContentClient.Close()
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	h := handler.New(renderer, pClient, ccClient, nClient)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux, ui.Files)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	return &App{
		cfg:                   cfg,
		profileClient:         pClient,
		notificationClient:    nClient,
		courseCatalogueClient: ccClient,
		courseContentClient:   cContentClient,
		server:                server,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Frontend Web Service running on port :%s...", a.cfg.Port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
	a.profileClient.Close()
	a.notificationClient.Close()
	a.courseCatalogueClient.Close()
}
