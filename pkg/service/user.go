package service

import (
	"dynamic-user-segmentation/pkg/model"
	"dynamic-user-segmentation/pkg/repository"
)

type UserService struct {
	repos repository.User
}

func NewUserService(repos repository.User) *UserService {
	return &UserService{repos: repos}
}

func (u *UserService) GetSegments(userId string) (model.UserSegments, error) {
	segments, err := u.repos.GetSegments(userId)
	if err != nil {
		return segments, err
	}

	return segments, nil
}

func (u *UserService) GetSegmentsHistory() ([]model.MemberEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, slugsToDelete []string) error {
	err := u.repos.UpdateSegments(userId, segmentsToAdd, slugsToDelete)
	if err != nil {
		return err
	}

	return nil
}
