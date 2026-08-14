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
	AssignmentServiceAddr      string
}

type App struct {
	cfg     Config
	clients *grpcclient.Clients
	server  *http.Server
}

func NewApp(cfg Config) (*App, error) {
	clients, err := grpcclient.Dial(
		cfg.ProfileServiceAddr,
		cfg.NotificationServiceAddr,
		cfg.CourseCatalogueServiceAddr,
		cfg.CourseContentServiceAddr,
		cfg.AssignmentServiceAddr,
	)
	if err != nil {
		return nil, err
	}

	renderer, err := render.New(ui.Files)
	if err != nil {
		clients.Close()
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	h := handler.New(renderer, clients)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux, ui.Files)

	return &App{
		cfg:     cfg,
		clients: clients,
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: mux,
		},
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
	a.clients.Close()
}
