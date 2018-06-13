package controllers

import (
	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type UsersController struct {
	Ctx iris.Context
	Service services.UserService
	Session *sessions.Session
}

type LoginReq struct {
	Name     string
	Password string
}

type CommonRes struct {
	State string
	Data  string
}

func (c *UsersController) PostLogin() (result CommonRes) {
	req := LoginReq{}
	c.Ctx.ReadForm(&req)
	valid, data, err := c.Service.Login(req.Name, req.Password)
	if err != nil { // 与Violet连接发生错误
		result.State = "error"
		result.Data = err.Error()
		return
	}
	if !valid { // 用户邮箱未激活
		result.State = "not_valid"
		result.Data = data
		return
	}

	userID, nikeName, tErr := c.Service.GetUser(data)
	if tErr != nil { // 无法获取用户详情
		result.State = "error"
		result.Data = tErr.Error()
		return
	}
	c.Session.Set("id", userID)

	result.State = "success"
	result.Data = nikeName
	return
}

type RegisterReq struct {
	Name     string
	Email    string
	Password string
}

func (c *UsersController) PostRegister() (res CommonRes) {
	req := RegisterReq{}
	c.Ctx.ReadForm(&req)
	err := c.Service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
	} else {
		res.State = "success"
	}
	return
}

func (c *UsersController) PostEmail() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.Data = "not_login"
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.Data = err.Error()
		return
	}
	err = c.Service.GetEmailCode(user.Email)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
	} else {
		res.State = "success"
	}
	return
}

type ValidReq struct {
	VCode string
}

func (c *UsersController) PostValid() (res CommonRes) {
	req := ValidReq{}
	res.State = "error"
	c.Ctx.ReadForm(&req)
	if c.Session.Get("id") == nil {
		res.Data = "not_login"
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.Data = err.Error()
		return
	}
	err = c.Service.ValidEmail(user.Email, req.VCode)
	if err != nil {
		res.Data = err.Error()
	} else {
		res.State = "success"
	}
	return
}

func (c *UsersController) PostLogout() (res CommonRes) {
	c.Session.Clear()
	res.State = "success"
	return
}

type SetNameReq struct {
	Name string
}
func (c *UsersController) PostUserName() (res CommonRes) {
	req := SetNameReq{}
	c.Ctx.ReadForm(&req)
	if c.Session.Get("id") == nil {
		res.State = "error"
		res.Data = "not_login"
		return
	}
	err := c.Service.SetUserName(c.Session.GetString("id"), req.Name)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
	} else {
		res.State = "success"
	}
	return
}

type UserRes struct {
	State string
	NikeName string
	Avatar string
	Gender int
	Level int
}
func (c *UsersController) GetUserBaseInfo() (res UserRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.State = "error"
		return
	}
	res.State = "success"
	res.NikeName = user.Info.NikeName
	res.Avatar = user.Info.Avatar
	res.Gender = user.Info.Gender
	res.Level = user.Level
	return
}