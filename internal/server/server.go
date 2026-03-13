package server

import (
	"context"
	"errors"
	"sync"

	"github.com/FulgurCode/BlankMap/config"
	// "github.com/FulgurCode/SevaNear/internal/db"
	// "github.com/FulgurCode/SevaNear/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

type WebServer struct {
	*fiber.App
	// DB     *db.Store
	Config *config.Config
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	var errs []error

	wg.Go(func() {
		if err := s.ShutdownWithContext(ctx); err != nil {
			errs = append(errs, err)
		}
	})

	// if s.DB != nil {
	// 	wg.Go(func() {
	// 		s.DB.Close(ctx)
	// 	})
	// }

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// func (s *WebServer) SetupMiddleware() {
// 	s.App.Use(middleware.SetupCORS())

// 	s.App.Use(middleware.SetupSession())
// }

func New(cfg *config.Config) *WebServer {
	// db := db.ConnectDB(cfg.DBString, cfg.MaxDBConns)

	var server = &WebServer{
		App: fiber.New(fiber.Config{
			AppName: "BlankMap - API",
		}),
		// DB:     db,
		Config: cfg,
	}

	// server.SetupMiddleware()

	return server
}
