package services

import "github.com/ZhenlyChen/BugServer/httpServer/models"

type Service struct {
	Model *models.Model
}

func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	return service
}

func (s *Service) NewUserService() UserService {
	return &userService{
		Model: s.Model.User,
	}
}
