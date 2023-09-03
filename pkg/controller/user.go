package controller

import (
	"dynamic-user-segmentation/pkg/model"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"strconv"
)

type UpdateSegmentsRequest struct {
	SegmentsToAdd []model.SegmentToAdd `json:"segments_to_add,omitempty"`
	SlugsToDelete []string             `json:"slugs_to_delete,omitempty"`
}

type LinkResponse struct {
	Link string `json:"link"`
}

func (c *Controller) getUserSegments(w http.ResponseWriter, r *http.Request) {
	log := c.log
	userId := chi.URLParam(r, "user_id")
	log.Info(fmt.Sprintf("getting segments for userId '%s", userId))
	segments, err := c.services.GetSegments(userId)
	if err != nil {
		log.Info(fmt.Sprintf("%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("User with id '%s' got active segments", userId)
	render.JSON(w, r, segments)
}

func (c *Controller) getSegmentsHistory(w http.ResponseWriter, r *http.Request) {
	log := c.log
	// получаем год и месяц
	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(fmt.Sprintf("Could not parse year to int got '%s'", yearStr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Error(fmt.Sprintf("Could not parse month to int got '%s'", monthStr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// получаем события за месяц
	fileName, err := c.services.GetSegmentsHistory(year, month)
	if err != nil {
		log.Error(fmt.Sprintf("Getting segments failed: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	link := "localhost:8082/api/public/v1/history/files/" + fileName
	resp := LinkResponse{Link: link}

	render.JSON(w, r, resp)
}

func (c *Controller) updateUserSegments(w http.ResponseWriter, r *http.Request) {
	log := c.log
	userId := chi.URLParam(r, "user_id")
	log.Info(fmt.Sprintf("updating segments for userId '%s'", userId))

	var request UpdateSegmentsRequest
	err := render.DecodeJSON(r.Body, &request)
	if errors.Is(err, io.EOF) {
		// Такую ошибку встретим, если получили запрос с пустым телом.
		// Обработаем её отдельно
		log.Error("request body is empty")

		render.JSON(w, r, "empty request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.services.UpdateSegments(userId, request.SegmentsToAdd, request.SlugsToDelete)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
