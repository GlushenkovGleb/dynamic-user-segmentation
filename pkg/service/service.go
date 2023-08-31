package service

import (
	"dynamic-user-segmentation/pkg/model"
	"dynamic-user-segmentation/pkg/repository"
)

type User interface {
	GetSegments(userId string) (model.UserSegments, error)
	GetSegmentsHistory() ([]model.MemberEvent, error)
	UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, segmentSlugsToDelete []string) error
}

type Segment interface {
	Create(segment model.Segment) (int, error)
	Delete(slug string) (int, error)
}

type Service struct {
	User
	Segment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repos),
		Segment: NewSegmentService(repos),
	}
}
