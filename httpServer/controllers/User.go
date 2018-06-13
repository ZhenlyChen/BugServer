package controllers

import (
	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type UsersController struct {
	// Optionally: context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserService

	Session *sessions.Session
}

func (c *UsersController) BeforeActivation(b mvc.BeforeActivation) {
	// b.Handle("GET", "/login", "Login")
}

func (c *UsersController) AfterActivation(a mvc.AfterActivation) {
	// fmt.Println(c.Session.Get("abc"))
}

type LoginReq struct {
	Name     string
	Password string
}

type CommonRes struct {
	State string
	Data string
}

func (c *UsersController) PostLogin() (result CommonRes) {
	req := LoginReq{}
	c.Ctx.ReadForm(&req)
	valid, data, err := c.Service.Login(req.Name, req.Password)
	if err != nil { // 与Violet连接发生错误
		result.State = "error"
		result.Data =  err.Error()
		return
	}
	if !valid { // 用户邮箱未激活
		result.State = "not_valid"
		result.Data = data
		return
	}

	userID, tErr := c.Service.GetUser(data)
	if tErr != nil { // 无法获取用户详情
		result.State = "error"
		result.Data = tErr.Error()
		return
	}
	c.Session.Set("id", userID)

	result.State = "success"
	return
}

type RegisterReq struct {
	Name     string
	Email 	 string
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
	email, err := c.Service.GetUserEmail(c.Session.GetString("id"))
	if err != nil {
		res.Data = err.Error()
		return
	}
	err = c.Service.GetEmailCode(email)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
	} else {
		res.State = "success"
	}
	return
}


type ValidReq struct {
	VCode 	 string
}

func (c *UsersController) PostValid() (res CommonRes) {
	req := ValidReq{}
	res.State = "error"
	c.Ctx.ReadForm(&req)
	if c.Session.Get("id") == nil {
		res.Data = "not_login"
		return
	}
	email, err := c.Service.GetUserEmail(c.Session.GetString("id"))
	if err != nil {
		res.Data = err.Error()
		return
	}
	err = c.Service.ValidEmail(email, req.VCode)
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