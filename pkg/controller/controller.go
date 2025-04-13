package controller

import (
	"log/slog"

	"dynamic-user-segmentation/pkg/repository"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo *repository.Repository
	log  *slog.Logger
}

func NewController(repo *repository.Repository, log *slog.Logger) *Controller {
	return &Controller{repo: repo, log: log}
}

func (c *Controller) InitRoutes(r *chi.Mux) {
	// группы
	r.Post("/groups", c.CreateGroup)
	r.Get("/groups", c.GetGroups)
	r.Get("/groups/{group_id}", c.GetGroup)
	r.Put("/groups/{group_id}", c.UpdateGroup)
	r.Delete("/groups/{group_id}", c.DeleteGroup)

	// студенты
	r.Post("/students", c.CreateStudent)
	r.Get("/students", c.GetStudents)
	r.Get("/students/{student_id}", c.GetStudent)
	r.Put("/students", c.UpdateStudent)
	r.Delete("/students/{student_id}", c.DeleteStudent)
}
