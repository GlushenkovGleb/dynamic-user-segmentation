package service

import (
	"dynamic-user-segmentation/pkg/model"
	"dynamic-user-segmentation/pkg/repository"
	"io"
)

type User interface {
	GetSegments(userId string) (model.UserSegments, error)
	GetSegmentsHistory(year, month int) (string, error)
	UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, segmentSlugsToDelete []string) error
}

type Segment interface {
	Create(segment model.Segment) (int, error)
	Delete(slug string) (int, error)
}

type History interface {
	UploadEvents(fileName string, dest io.Writer) (int, error)
}

type Service struct {
	User
	Segment
	History
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repos.User, repos.History),
		Segment: NewSegmentService(repos),
		History: NewHistoryService(repos),
	}
}
