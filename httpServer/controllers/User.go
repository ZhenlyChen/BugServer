package controllers

import (
	"html/template"
	"regexp"

	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// UsersController 用户控制
type UsersController struct {
	Ctx     iris.Context
	Service services.UserService
	Session *sessions.Session
}

// LoginReq OST /user/login 登陆请求
type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// PostLogin POST /user/login 登陆
func (c *UsersController) PostLogin() (result CommonRes) {
	req := LoginReq{}
	c.Ctx.ReadJSON(&req)
	if req.Name == "" || req.Password == "" {
		result.Status = StatusBadReq
		return
	}
	valid, data, err := c.Service.Login(req.Name, req.Password)
	if err != nil { // 与Violet连接发生错误
		result.Status = StatusError
		result.Msg = err.Error()
		return
	}
	if !valid { // 用户邮箱未激活
		result.Status = StatusNotValid
		result.Msg = data
		return
	}

	userID, nikeName, tErr := c.Service.GetUserFromViolet(data)
	if tErr != nil { // 无法获取用户详情
		result.Status = StatusError
		result.Msg = tErr.Error()
		return
	}
	c.Session.Set("id", userID)

	result.Status = StatusSuccess
	result.Msg = nikeName
	return
}

// RegisterReq POST /user/register 注册请求
type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PostRegister POST /user/register 注册
func (c *UsersController) PostRegister() (res CommonRes) {
	req := RegisterReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.Status = StatusBadReq
	}
	if err := c.Service.Register(req.Name, req.Email, req.Password); err != nil {
		res.Status = err.Error()
	} else {
		res.Status = StatusSuccess
	}
	return
}

// PostEmail POST /user/email 获取邮箱验证码
func (c *UsersController) PostEmail() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.Status = err.Error()
		return
	}
	if err := c.Service.GetEmailCode(user.Email); err != nil {
		res.Status = err.Error()
	} else {
		res.Status = StatusSuccess
	}
	return
}

// ValidReq POST /user/valid/ 请求
type ValidReq struct {
	VCode string `json:"vCode"`
}

// PostValid POST /user/valid/ 验证邮箱
func (c *UsersController) PostValid() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	req := ValidReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.Status = StatusBadReq
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.Msg = err.Error()
		return
	}
	if err := c.Service.ValidEmail(user.Email, req.VCode); err != nil {
		res.Status = err.Error()
	} else {
		res.Status = StatusSuccess
	}
	return
}

// PostLogout POST /user/logout 退出登陆
func (c *UsersController) PostLogout() (res CommonRes) {
	c.Session.Clear()
	res.Status = StatusSuccess
	return
}

// InfoReq POST /user/info 请求结构
type InfoReq struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}

// PostInfo POST /user/info 更新信息
func (c *UsersController) PostInfo() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 检测姓名合法性
	req := InfoReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil || req.Name == "" || len(req.Name) > 20 {
		res.Status = StatusBadReq
		return
	}
	// 检测非法字符
	if m, _ := regexp.MatchString(`[\\\/\(\)<|> "'{}:;]`, req.Name); m {
		res.Status = StatusBadReq
		return
	}
	req.Name = template.HTMLEscapeString(req.Name)
	if err := c.Service.SetUserInfo(c.Session.GetString("id"), models.UserInfo{
		NikeName: req.Name,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
	}); err != nil {
		res.Status = err.Error()
	} else {
		res.Status = StatusSuccess
	}
	return
}

// UserRes GET /user/info/{userID} 返回值
type UserRes struct {
	Status   string `json:"status"`
	NikeName string `json:"nikeName"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Level    int    `json:"level"`
}

// GetInfoBy GET /user/info/{userID} 获取用户信息，userID为空时候获取自身信息
func (c *UsersController) GetInfoBy(id string) (res UserRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	if id != "" && !bson.IsObjectIdHex(id) {
		res.Status = StatusBadReq
		return
	} else if id == "" {
		id = c.Session.GetString("id")
	}
	user, err := c.Service.GetUserInfo(id)
	if err != nil {
		res.Status = StatusError
		return
	}
	res.Status = StatusSuccess
	res.NikeName = user.Info.NikeName
	res.Avatar = user.Info.Avatar
	res.Gender = user.Info.Gender
	res.Level = user.Level
	return
}
