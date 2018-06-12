package services

import (
	violet "github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/kataras/iris/core/errors"
	"encoding/json"
	"fmt"
)

type UserService interface {
	Login(name, password string) (valid bool, email string, err error)
	GetUser(code string) (ID string, err error)
	InitViolet(c violet.Config)
	// Register(name, email, password string) error
	// GetEmailCode(email string) error
}

type userService struct {
	Model models.UserModel
	Violet violet.Violet
}

type LoginRes struct {
	Valid bool
	Email string
	Code  string
}

func (s *userService) Login(name, password string) (valid bool, data string, err error) {
	resp, tErr := s.Violet.Login(name, password)
	if tErr != nil {
		err = tErr
		return
	}
	// 非正常的返回码
	if resp.StatusCode() != 200 {
		err = errors.New(resp.String())
		return
	}
	// 解析结果
	var loginRes LoginRes
	err = json.Unmarshal([]byte(resp.String()), &loginRes)
	if err != nil {
		return
	}
	valid = loginRes.Valid
	// 未激活邮箱
	if !valid {
		data = loginRes.Email
		return
	}
	// 登陆成功
	valid = true
	data = loginRes.Code
	return

}

type TokenRes struct {
	UserID string
	Token string
}

func (s *userService) GetUser(code string) (ID string, err error) {
	// 获取用户Token
	resp, tErr := s.Violet.GetToken(code)
	if tErr != nil {
		err = tErr
		return
	}
	if resp.StatusCode() != 200 {
		err = errors.New(resp.String())
		return
	}
	// 解析结果
	var tokenRes TokenRes
	err = json.Unmarshal([]byte(resp.String()), &tokenRes)
	ID = tokenRes.UserID
	fmt.Println(tokenRes)
	// 保存数据并获取用户信息
	// resp, err = s.Violet.GetUserBaseInfo()

	return
}



func (s *userService) InitViolet(c violet.Config) {
	s.Violet = violet.NewViolet(c)
}
