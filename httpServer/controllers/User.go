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

type LoginRes struct {
	State string
	Data string
}

func (c *UsersController) PostLogin() (result LoginRes) {
	req := LoginReq{}
	c.Ctx.ReadForm(&req)
	valid, data, err := c.Service.Login(req.Name, req.Password)
	if err != nil {
		result.State = "error"
		result.Data =  err.Error()
		return
	}
	if !valid {
		result.State = "not_valid"
		result.Data = data
		return
	}

	userID, tErr := c.Service.GetUser(data)
	if tErr != nil {
		result.State = "error"
		result.Data = tErr.Error()
		return
	}
	c.Session.Set("id", userID)

	result.State = "success"
	return
}
