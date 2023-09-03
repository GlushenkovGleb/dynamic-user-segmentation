package service

import (
	"dynamic-user-segmentation/pkg/repository"
	"io"
)

type HistoryService struct {
	repos repository.History
}

func NewHistoryService(repos repository.History) *HistoryService {
	return &HistoryService{repos: repos}
}

// UploadEvents takes fileNane and dest to write,
// returns fileSize
func (h *HistoryService) UploadEvents(fileName string, dest io.Writer) (int, error) {
	fileSize, err := h.repos.UploadEvents(fileName, dest)
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}
