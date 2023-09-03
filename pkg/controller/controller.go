package controller

import (
	"dynamic-user-segmentation/pkg/service"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Controller struct {
	services *service.Service
	log      *slog.Logger
}

func NewController(services *service.Service, log *slog.Logger) *Controller {
	return &Controller{services: services, log: log}
}

func (c *Controller) InitRoutes(r *chi.Mux) {
	r.Route("/api/public/v1", func(r chi.Router) {
		// пользователи
		r.Get("/users/{user_id}/segments", c.getUserSegments)
		r.Get("/users/history/{year}-{month}", c.getSegmentsHistory)
		r.Put("/users/{user_id}/segments", c.updateUserSegments)

		// сегменты
		r.Post("/segments", c.createSegment)
		r.Delete("/segments/{slug}", c.deleteSegment)

		// история
		r.Get("/history/files/{file_name}", c.uploadHistory)
	})
}
