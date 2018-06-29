package services

import "github.com/ZhenlyChen/BugServer/httpServer/models"

type Service struct {
	Model *models.Model
	User userService
}

func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	service.User = userService{
		Model: &m.User,
		Service: service,
		UserInfo: make(map[string]UserBaseInfo),
	}
	return service
}

func (s *Service) NewUserService() UserService {
	return &s.User
}
