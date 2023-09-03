package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (c *Controller) uploadHistory(w http.ResponseWriter, r *http.Request) {
	log := c.log
	fileName := chi.URLParam(r, "file_name") + ".csv"
	fileContentType := "text/csv"

	w.Header().Set("Content-Type", fileContentType+";"+fileName)
	_, err := c.services.UploadEvents(fileName, w)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
