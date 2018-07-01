package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"strconv"
	"github.com/globalsign/mgo/bson"
)

// UsersController 用户控制
type RoomsController struct {
	Ctx     iris.Context
	Service services.RoomService
	Session *sessions.Session
}

type RoomsRes struct {
	Status string              `json:"status"`
	Count  int                 `json:"count"`
	Rooms  []services.GameRoom `json:"rooms"`
}

// GetRooms GET /room/list/{page} 获取房间列表（每页10个）page:1~10
func (c *RoomsController) GetListBy(pageStr string) (res RoomsRes) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		res.Status = err.Error()
		return
	}
	if page < 1 || page > 10 {
		res.Status = StatusBadReq
		return
	}

	rooms := c.Service.GetRooms()
	res.Count = len(rooms)
	// 删除密码
	for _, room := range rooms {
		if room.Password != "" {
			room.Password = "password"
		}
	}
	endIndex := page * 10
	if res.Count < (page-1)*10 {
		res.Status = StatusNull
		return
	} else if res.Count < page*10 {
		endIndex = res.Count - 1
	}

	res.Status = StatusSuccess
	res.Rooms = rooms[(page-1)*10 : endIndex]
	return
}

type RoomRes struct {
	Status     string                `json:"status"`
	Room       services.GameRoom     `json:"room"`
	PlayerInfo []services.PlayerInfo `json:"players"`
}

// GetRoom GET /room/detail/{roomID} 获取单个房间详情
func (c *RoomsController) GetDetailBy(id string) (res RoomRes) {
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
	room, err := c.Service.GetRoom(roomID)
	if err != nil {
		res.Status = err.Error()
		return
	}
	// 是否在房间内
	inRoom := false
	userID := c.Session.GetString("id")
	for _, player := range room.Players {
		if player.UserID == userID {
			inRoom = true
		}
	}
	if !inRoom {
		// 不再房间内
		res.Status = StatusNotAllow
		return
	}
	// 获取玩家详细信息
	playInfo, err := c.Service.GetPlayers(roomID)
	if err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	res.Room = *room
	res.PlayerInfo = playInfo
	if res.Room.Password != "" {
		res.Room.Password = "password"
	}
	return
}

type reqNewRoom struct {
	Title     string `json:"title"`
	Password  string `json:"password"`
	GameMap   string `json:"gameMap"`
	GameMode  string `json:"gameMode"`
	MaxPlayer int    `josn:"maxPlayer"`
}

// PostRoom POST /room/new 新建并加入房间
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
	roomID, err := c.Service.AddRoom(c.Session.GetString("id"), req.Title, req.GameMode, req.GameMap, req.Password, req.MaxPlayer)
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

// PostJoin POST /room/join/{roomId} 加入房间
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

// PostReady POST /room/ready/{true/false} 设置准备状态
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

// PostTeam POST /room/team/{teamID} 设置队伍
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

// PostRole POST /room/role/{roleName} 设置角色
func (c *RoomsController) PostRoleBy(role string) (res CommonRes) {
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
	if role == "" {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetRole(roomID, c.Session.GetString("id"), role); err != nil {
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
	if err := c.Service.QuitRoom(roomID, c.Session.GetString("id")); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	c.Session.Delete("room")
	return
}

type roomInfoReq struct {
	MaxPlayer int    `json:"maxPlayer"`
	GameMap   string `json:"gameMap"`
	GameMode  string `json:"gameMode"`
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
	req := roomInfoReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.Status = StatusBadReq
		return
	}
	if err := c.Service.SetRoomInfo(roomID, req.MaxPlayer, c.Session.GetString("id"), req.GameMap, req.GameMode); err != nil {
		res.Status = err.Error()
		return
	}
	res.Status = StatusSuccess
	return
}

// PostRole POST /room/own/{userId} 设置房主
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

// PostRole POST /room/out/{userId} 踢人
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
