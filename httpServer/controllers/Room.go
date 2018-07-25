package controllers

import (
	"strconv"

	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// RoomsController 用户控制
type RoomsController struct {
	Ctx     iris.Context
	Service services.RoomService
	Session *sessions.Session
}

// RoomsRes 房间列表
type RoomsRes struct {
	Status string              `json:"status"`
	Count  int                 `json:"count"`
	Rooms  []services.RoomInfo `json:"rooms"`
}

// GetHeart GET /room/heart 发送心跳包
func (c *RoomsController) GetHeart() string {
	if c.Session.Get("id") == nil || c.Session.Get("room") == nil {
		return "false"
	}
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		return "false"
	}
	if c.Service.Heart(c.Session.GetString("id"), roomID) {
		return "true"
	}
	return "false"
}

// GetListBy GET /room/list/{page}?size=n 获取房间列表（每页n(1-20)个）
func (c *RoomsController) GetListBy(pageStr string) (res RoomsRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		res.Status = StatusBadReq
		return
	}
	pageSize, err := strconv.Atoi(c.Ctx.FormValue("size"))
	if err != nil {
		res.Status = StatusBadReq
		return
	}
	if pageSize < 1 || pageSize > 20 {
		res.Status = StatusBadReq
		return
	}
	if page < 1 || page > pageSize {
		res.Status = StatusBadReq
		return
	}
	endIndex := page * pageSize
	if res.Count <= (page-1)*pageSize {
		res.Status = StatusNull
		return
	} else if res.Count < page*pageSize {
		endIndex = res.Count
	}
	rooms := c.Service.GetRooms()
	res.Count = len(rooms)
	// 删除密码 处理用户
	for i := (page - 1) * pageSize; i < endIndex; i++ {
		room := &rooms[i]
		if room.Info.Password != "" {
			room.Info.Password = "password"
		}
		room.Info.PlayerCount = len(room.Info.Players)
		room.Info.Players = []services.Player{}
		res.Rooms = append(res.Rooms, room.Info)
	}
	res.Status = StatusSuccess
	return
}

// RoomRes ...
type RoomRes struct {
	Status   string            `json:"status"`
	RoomInfo services.RoomInfo `json:"room"`
}

// GetDetail GET /room/detail/ 获取自己房间详情
func (c *RoomsController) GetDetail() (res RoomRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	room, err := c.Service.GetRoom(roomID)
	if err != nil {
		res.Status = err.Error()
		return
	}
	// 是否已经被踢出房间
	userID := c.Session.GetString("id")
	if !c.Service.IsInRoom(userID) {
		res.Status = StatusNotFound
		return
	}
	// 获取玩家详细信息
	res.Status = StatusSuccess
	res.RoomInfo = (*room).Info
	if res.RoomInfo.Password != "" {
		res.RoomInfo.Password = "password"
	}
	return
}

type reqNewRoom struct {
	Title     string `json:"title"`
	Password  string `json:"password"`
	GameMap   string `json:"gameMap"`
	GameMode  string `json:"gameMode"`
	MaxPlayer int    `json:"maxPlayer"`
	IsRandom  bool   `json:"isRandom"`
}

// PostNew POST /room/new 新建并加入房间
func (c *RoomsController) PostNew() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	req := reqNewRoom{}
	if err := c.Ctx.ReadJSON(&req); err != nil ||
		req.GameMap == "" || req.Title == "" ||
		req.GameMode == "" || req.MaxPlayer < 0 ||
		req.MaxPlayer > services.MaxPlayer {
		res.Status = StatusBadReq
		return
	}
	roomID, err := c.Service.AddRoom(c.Session.GetString("id"), req.Title, req.GameMode, req.GameMap, req.Password, req.MaxPlayer, req.IsRandom)
	if err != nil {
		res.Status = err.Error()
		return
	}
	c.Session.Set("room", roomID)
	res.Status = StatusSuccess
	res.Msg = strconv.Itoa(roomID)
	return
}

type reqJoinRoom struct {
	Password string `json:"password"`
}

// PostJoinBy POST /room/join/{roomId} 加入房间
func (c *RoomsController) PostJoinBy(id string) (res CommonRes) {
	// 检测参数合法性
	roomID, err := strconv.Atoi(id)
	if err != nil {
		res.Status = StatusBadReq
		return
	}
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	req := reqJoinRoom{}
	password := ""
	if err := c.Ctx.ReadJSON(&req); err == nil {
		password = req.Password
	}
	if err := c.Service.JoinRoom(roomID, c.Session.GetString("id"), password); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	c.Session.Set("room", roomID)
	return
}

// PostReadyBy POST /room/ready/{true/false} 设置准备状态
func (c *RoomsController) PostReadyBy(isReady string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if isReady != "true" && isReady != "false" {
		res.Status = StatusBadReq
		return
	}
	ready := false
	if isReady == "true" {
		ready = true
	}
	if err := c.Service.SetReady(roomID, c.Session.GetString("id"), ready); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostReadyBy POST /room/play/{true/false} 设置开始状态
func (c *RoomsController) PostPlayBy(isPlaying string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if isPlaying != "true" && isPlaying != "false" {
		res.Status = StatusBadReq
		return
	}
	playing := false
	if isPlaying == "true" {
		playing = true
	}
	if err := c.Service.SetPlaying(roomID, c.Session.GetString("id"), playing); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostTeamBy POST /room/team/{teamID} 设置队伍
func (c *RoomsController) PostTeamBy(teamStr string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	teamID, err := strconv.Atoi(teamStr)
	if err != nil {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetTeam(roomID, teamID, c.Session.GetString("id")); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostRoleBy POST /room/role/{roleName} 设置角色
func (c *RoomsController) PostRoleBy(roleStr string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if roleStr == "" {
		res.Status = StatusBadReq
		return
	}
	roleID, err := strconv.Atoi(roleStr)
	if err != nil {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetRole(roomID, roleID, c.Session.GetString("id")); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostQuit POST /room/quit 退出房间
func (c *RoomsController) PostQuit() (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	room, err := c.Service.GetRoom(roomID)
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if err := c.Service.QuitRoom(room, c.Session.GetString("id")); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	c.Session.Delete("room")
	return
}

// PostInfo POST /room/info 设置房间信息
func (c *RoomsController) PostInfo() (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	req := reqNewRoom{}
	if err := c.Ctx.ReadJSON(&req); err != nil ||
		req.GameMap == "" || req.Title == "" || req.MaxPlayer < 0 ||
		req.MaxPlayer > services.MaxPlayer {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetRoomInfo(roomID, req.MaxPlayer, c.Session.GetString("id"), req.GameMap, req.Title, req.Password, req.IsRandom); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostOwnBy POST /room/own/{userId} 设置房主
func (c *RoomsController) PostOwnBy(id string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetRoomOwn(roomID, c.Session.GetString("id"), id); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostOutBy POST /room/out/{userId} 踢人
func (c *RoomsController) PostOutBy(id string) (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.GetOutRoom(roomID, c.Session.GetString("id"), id); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostStart POST /room/start 开始游戏
func (c *RoomsController) PostStart() (res CommonRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	if err := c.Service.StartGame(roomID, c.Session.GetString("id")); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

type KeyRes struct {
	Status string `json:"status"`
	Key    int    `json:"key"`
}

// GetKey GET /room/key 获取进入对局的密钥
func (c *RoomsController) GetKey() (res KeyRes) {
	// 是否登陆
	if c.Session.Get("id") == nil {
		res.Status = StatusNotLogin
		return
	}
	// 是否在房间里面
	roomID, err := c.Session.GetInt("room")
	if err != nil {
		res.Status = StatusNotFound
		return
	}
	key, err := c.Service.GetKey(c.Session.GetString("id"), roomID)
	if err != nil {
		res.Status = StatusNotAllow
		return
	}
	res.Key = key
	res.Status = StatusSuccess
	return
}
