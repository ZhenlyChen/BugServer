package services

import (
	"fmt"

	violet "github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
)

type UserService interface {
	Login(name, password string) (valid bool, email string, err error)
	// Register(name, email, password string) error
	// GetEmailCode(email string) error
}

type userService struct {
	Model models.UserModel
}

func (s *userService) Login(name, password string) (valid bool, email string, err error) {
	res, tErr := violet.Login(name, password)
	fmt.Println(res)
	fmt.Println(tErr)
	if tErr != nil {
		err = tErr
		return
	}
	return
}
