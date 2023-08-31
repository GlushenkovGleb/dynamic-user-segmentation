package service

import (
	"dynamic-user-segmentation/pkg/model"
	"dynamic-user-segmentation/pkg/repository"
)

type SegmentService struct {
	repos repository.Segment
}

func NewSegmentService(repos repository.Segment) *SegmentService {
	return &SegmentService{repos: repos}
}

func (s *SegmentService) Create(segment model.Segment) (int, error) {
	segId, err := s.repos.CreateSegment(segment)
	if err != nil {
		return 0, err
	}
	return segId, err
}

func (s *SegmentService) Delete(slug string) (int, error) {
	id, err := s.repos.DeleteSegment(slug)
	if err != nil {
		return 0, err
	}

	return id, nil
}
