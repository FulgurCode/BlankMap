package server

import (
	"github.com/FulgurCode/BlankMap/internal/handlers"
	"github.com/FulgurCode/BlankMap/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func (s *WebServer) RegisterRoutes() {
	authHandler := handlers.NewAuthHandler(s.DB.Queries, s.Config.JWTSecret)
	blankMapHandler := handlers.NewBlankMapHandler(s.DB.Queries)
	pinHandler := handlers.NewPinHandler(s.DB.Queries)
	// feedbackHandler := handlers.NewFeedbackHandler(s.DB.Queries)

	s.App.Get("/", s.handleIndex)

	// Auth — public
	auth := s.App.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	api := s.App.Group("/", middleware.RequireAuth(s.Config.JWTSecret))

	api.Get("/auth/me", authHandler.Me)

	// Blank Maps
	api.Post("/blank-maps", blankMapHandler.CreateBlankMap)
	api.Get("/blank-maps", blankMapHandler.GetBlankMaps)
	api.Get("/blank-maps/:id", blankMapHandler.GetBlankMapByID)
	api.Put("/blank-maps/:id", blankMapHandler.UpdateBlankMap)
	api.Delete("/blank-maps/:id", blankMapHandler.DeleteBlankMap)

	// Pins
	// GET /pins?blank_map_id=<uuid>  — filter by map
	// GET /pins/nearby?lat=&lng=&radius=  — proximity search
	api.Post("/pins", pinHandler.CreatePin)
	api.Get("/pins", pinHandler.GetPins)
	api.Get("/pins/nearby", pinHandler.GetPinsNearby)
	api.Get("/pins/:id", pinHandler.GetPinByID)
	api.Put("/pins/:id", pinHandler.UpdatePin)
	api.Delete("/pins/:id", pinHandler.DeletePin)

	// // Feedback
	// api.Post("/feedback", feedbackHandler.CreateFeedback)
	// api.Get("/pins/:pinID/feedback", feedbackHandler.GetFeedbackByPin)
	// api.Put("/feedback/:id", feedbackHandler.UpdateFeedback)
	// api.Delete("/feedback/:id", feedbackHandler.DeleteFeedback)
}

func (s *WebServer) handleIndex(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"app":     "BlankMap API",
		"version": "1.0.0",
	})
}
