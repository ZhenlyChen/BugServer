package services

import (
	"encoding/json"

	violet "github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/kataras/iris/core/errors"
)

type UserService interface {
	InitViolet(c violet.Config)
	// 登陆部分API
	Login(name, password string) (valid bool, email string, err error)
	GetUser(code string) (ID string, err error)
	Register(name, email, password string) (err error)
	GetEmailCode(email string) error
	ValidEmail(email, vCode string) error
	GetUserEmail(id string) (email string, err error)
}

type userService struct {
	Model  models.UserModel
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
	Token  string
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

	// 保存数据并获取用户信息
	user, tErr2 := s.Model.GetUserByVID(tokenRes.UserID)
	if tErr2 == nil { // 数据库已存在该用户
		ID = user.ID.Hex()
		s.Model.SetUserToken(user.ID.Hex(), tokenRes.Token)
	} else if tErr2.Error() == "not found" { // 数据库不存在此用户
		ID, err = s.SaveUser(tokenRes.UserID, tokenRes.Token)
	} else { // 其他错误
		err = tErr2
	}
	return
}

type UserInfoRes struct {
	Email string
	Name  string
	Info  UserInfo
}

type UserInfo struct {
	Avatar string
	Gender int
}

func (s *userService) SaveUser(userVID, token string) (ID string, err error) {
	resp, err := s.Violet.GetUserBaseInfo(userVID, token)
	if err != nil {
		return
	}
	if resp.StatusCode() != 200 {
		err = errors.New(resp.String())
		return
	}
	// 解析结果
	var userInfoRes UserInfoRes
	err = json.Unmarshal([]byte(resp.String()), &userInfoRes)
	if err != nil {
		return
	}

	bsonID, err := s.Model.AddUser(userVID, userInfoRes.Name, userInfoRes.Email, token, userInfoRes.Info.Avatar, userInfoRes.Info.Gender)
	ID = bsonID.Hex()
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

func (s *userService) GetUserEmail(id string) (email string, err error) {
	user, err := s.Model.GetUserByID(id)
	if err == nil {
		email = user.Email
	}
	return
}

func (s *userService) InitViolet(c violet.Config) {
	s.Violet = violet.NewViolet(c)
}
