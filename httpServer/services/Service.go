package services

import "github.com/ZhenlyChen/BugServer/httpServer/models"

type Service struct {
	Model *models.Model
	User userService
}

func NewSerivce(m *models.Model) *Service {
	service := new(Service)
	service.User.Model = m
	return service
}