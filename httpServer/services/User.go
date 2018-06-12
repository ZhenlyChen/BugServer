package services

import (
	"strings"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
)

type UserService interface {
	Login(name, password string) (vaild bool, email string, err error)
	Register(name, email, password string) error
	GetEmailCode(email string) error
}

type userService struct {
	Model *models.UserModel
}

func (s *userService) Login(name, password string) (vaild bool, email string, err error) {
	if (strings.Index(name, "@") != -1) {
		user := s.Model.
	}
	return
}
