package services

import "github.com/ZhenlyChen/BugServer/httpServer/models"

type UserService interface {

}

type userService struct {
	Model *models.Model
}