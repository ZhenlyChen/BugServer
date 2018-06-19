package services

import (
	"encoding/json"

	violet "github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/kataras/iris/core/errors"
)

// UserService 用户服务
type UserService interface {
	InitViolet(c violet.Config)
	// 登陆部分API
	Login(name, password string) (valid bool, email string, err error)
	GetUser(code string) (ID, name string, err error)
	Register(name, email, password string) (err error)
	GetEmailCode(email string) error
	ValidEmail(email, vCode string) error
	GetUserInfo(id string) (user models.Users, err error)
	SetUserName(id, name string) error
}

type userService struct {
	Model  models.UserModel
	Violet violet.Violet
}

type loginRes struct {
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
	var res loginRes
	err = json.Unmarshal([]byte(resp.String()), &res)
	if err != nil {
		return
	}
	valid = res.Valid
	// 未激活邮箱
	if !valid {
		data = res.Email
		return
	}
	// 登陆成功
	valid = true
	data = res.Code
	return

}

func (s *userService) GetUser(code string) (ID, name string, err error) {
	// 获取用户Token
	tokenRes, err := s.Violet.GetToken(code)
	if err != nil {
		return
	}
	// 保存数据并获取用户信息
	user, tErr2 := s.Model.GetUserByVID(tokenRes.UserID)
	if tErr2 == nil { // 数据库已存在该用户
		ID = user.ID.Hex()
		name = user.Info.NikeName
		s.Model.SetUserToken(user.ID.Hex(), tokenRes.Token)
	} else if tErr2.Error() == "not found" { // 数据库不存在此用户
		userInfoRes, err := s.Violet.GetUserBaseInfo(tokenRes.UserID, tokenRes.Token)
		if err != nil {
			return "", "", err
		}
		bsonID, err := s.Model.AddUser(tokenRes.UserID, userInfoRes.Name, userInfoRes.Email, tokenRes.Token, userInfoRes.Info.Avatar, userInfoRes.Info.Gender)
		ID = bsonID.Hex()
		name = "new_user"
	} else { // 其他错误
		err = tErr2
	}
	return
}

func (s *userService) Register(name, email, password string) error {
	resp, err := s.Violet.Register(name, email, password)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}
	return nil
}

func (s *userService) GetEmailCode(email string) error {
	resp, err := s.Violet.GetEmailCode(email)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}
	return nil
}

func (s *userService) ValidEmail(email, vCode string) error {
	resp, err := s.Violet.ValidEmail(email, vCode)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}
	return nil
}

func (s *userService) GetUserInfo(id string) (user models.Users, err error) {
	user, err = s.Model.GetUserByID(id)
	return
}

func (s *userService) SetUserName(id, name string) error {
	err := s.Model.SetUserName(id, name)
	return err
}

func (s *userService) InitViolet(c violet.Config) {
	s.Violet = violet.NewViolet(c)
}
