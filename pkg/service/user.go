package service

import (
	"dynamic-user-segmentation/pkg/model"
	"dynamic-user-segmentation/pkg/repository"
)

type UserService struct {
	userRep    repository.User
	historyRep repository.History
}

func NewUserService(userRep repository.User, historyRep repository.History) *UserService {
	return &UserService{userRep: userRep, historyRep: historyRep}
}

func (u *UserService) GetSegments(userId string) (model.UserSegments, error) {
	segments, err := u.userRep.GetSegments(userId)
	if err != nil {
		return segments, err
	}

	return segments, nil
}

func (u *UserService) GetSegmentsHistory(year, month int) (string, error) {
	events, err := u.userRep.GetSegmentsHistory(year, month)
	if err != nil {
		return "", err
	}

	fileName, err := u.historyRep.SaveEvents(events)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (u *UserService) UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, slugsToDelete []string) error {
	err := u.userRep.UpdateSegments(userId, segmentsToAdd, slugsToDelete)
	if err != nil {
		return err
	}
	return nil
}
