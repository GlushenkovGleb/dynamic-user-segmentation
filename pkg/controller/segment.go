package controller

import (
	"dynamic-user-segmentation/pkg/model"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

func (c *Controller) createSegment(w http.ResponseWriter, r *http.Request) {
	log := c.log
	var segment model.Segment

	err := render.DecodeJSON(r.Body, &segment)
	if errors.Is(err, io.EOF) {
		// Такую ошибку встретим, если получили запрос с пустым телом.
		// Обработаем её отдельно
		log.Error("request body is empty")

		render.JSON(w, r, "empty request")

		return
	}

	log.Info("request body decoded", slog.Any("request", segment))

	if err := validator.New().Struct(segment); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request")

		render.JSON(w, r, validateErr)

		return
	}

	id, err := c.services.Create(segment)
	if err != nil {
		log.Error(fmt.Sprintf("%s"), err)
		return
	}
	log.Info("Segment added!", id)

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) deleteSegment(w http.ResponseWriter, r *http.Request) {
	log := c.log
	slug := chi.URLParam(r, "slug")
	_, err := c.services.Delete(slug)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info(fmt.Sprintf("Deleted segment with slug '%s'", slug))

	w.WriteHeader(http.StatusNoContent)
}
