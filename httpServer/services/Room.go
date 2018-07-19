package services

import (
	"time"

	"github.com/ZhenlyChen/BugServer/gameServer"
)

// RoomService ...
type RoomService interface {
	InitGameServer(config gameServer.ServerConfig)
	CheckHeart()
	Heart(userID string, roomID int) bool
	GetRoom(roomID int) (room *GameRoom, err error)
	JoinRoom(roomID int, userID, password string) error
	SetReady(roomID int, userID string, isReady bool) error
	SetRole(roomID int, userID, role string) error
	SetTeam(roomID, teamID int, userID string) error
	GetRooms() []GameRoom
	GetPlayers(roomID int) (info []PlayerInfo, err error)
	QuitRoom(roomID int, userID string) error
	// 房主
	StartGame(roomID int, ownID string) error
	SetRoomOwn(roomID int, ownID, newOwnID string) error
	AddRoom(ownID, title, mode, gameMap, password string, maxPlayer int) (roomID int, err error)
	SetRoomInfo(roomID, maxPlayer int, ownID, gameMap, gameMode string) error
	GetOutRoom(roomID int, ownID, userID string) error
}

type roomService struct {
	Service *Service
	Game    *gameServer.GameServer
	Rooms   []GameRoom
}

// GameMode 游戏模式
const (
	GameModePersonal  = "personal" // 个人
	GameModeTogether  = "together" // 合作
	GameModeTeamTwo   = "team2"    // 2人团队
	GameModeTeamThree = "team3"    // 3人团队
	GameModeTeamFour  = "team4"    // 4人团队

	MaxRoom   = 100
	MaxPlayer = 20
)

// GameRoom 房间数据
type GameRoom struct {
	ID        int      `json:"id"`        // 房间 ID
	OwnID     string   `json:"ownId"`     // 房主ID
	Port      int      `json:"port"`      // 房间服务器端口
	Title     string   `json:"title"`     // 标题
	GameMap   string   `json:"gameMap"`   // 游戏地图
	MaxPlayer int      `json:"maxPlayer"` // 最大人数
	Mode      string   `json:"mode"`      // 游戏模式
	Password  string   `json:"password"`  // 房间密码
	Playing   bool     `json:"playing"`   // 是否正在玩
	Players   []Player `json:"players"`   // 玩家数据
}

// Player 玩家信息
type Player struct {
	UserID  string `json:"userId"`  // 玩家ID
	GameID  int    `json:"gameId"`  // 游戏内ID
	RoleID  string `json:"roleId"`  // 角色ID
	IsReady bool   `json:"isReady"` // 是否准备
	Heart   int    `json:"heart"`   // 心跳💗
	Team    int    `json:"team"`    // "1-4" - 队伍一~四
}

func (s *roomService) CheckHeart() {
	for {
		time.Sleep(time.Second)
		for i := range s.Rooms {
			for j := range s.Rooms[i].Players {
				if s.Rooms[i].Players[j].Heart > 3 {
					s.QuitRoom(s.Rooms[i].ID, s.Rooms[i].Players[j].UserID)
				} else {
					s.Rooms[i].Players[j].Heart++
				}
			}
		}
	}
}

func (s *roomService) Heart(userID string, roomID int) bool {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return false
	}
	for i := range room.Players {
		if room.Players[i].UserID == userID {
			room.Players[i].Heart = 0
			return true
		}
	}
	return false
}

func (s *roomService) InitGameServer(config gameServer.ServerConfig) {
	s.Game = new(gameServer.GameServer)
	s.Game.InitServer(config)
}

// GetRooms 获取房间列表
func (s *roomService) GetRooms() []GameRoom {
	return s.Rooms
}

// PlayerInfo 玩家个性信息
type PlayerInfo struct {
	Player Player       `json:"player"`
	Info   UserBaseInfo `json:"info"`
}

// findRoom 寻找房间，返回房间地址
func (s *roomService) GetRoom(roomID int) (room *GameRoom, err error) {
	for i := range s.Rooms {
		if roomID == s.Rooms[i].ID {
			room = &s.Rooms[i]
			return
		}
	}
	err = ErrNotFound
	return
}

// getGameID 获取新的游戏ID
func (s *roomService) newGameID(players []Player) int {
	for i := 0; i < MaxPlayer; i++ {
		hasExist := false
		for _, player := range players {
			if player.GameID == i {
				hasExist = true
				break
			}
		}
		if !hasExist {
			return i
		}
	}
	return -1
}

// newRoomID 获取新的房间ID
func (s *roomService) newRoomID() int {
	for i := 0; i < MaxRoom; i++ {
		hasExist := false
		for _, room := range s.Rooms {
			if room.ID == i {
				hasExist = true
				break
			}
		}
		if !hasExist {
			return i
		}
	}
	return -1
}

func (s *roomService) getTeamID(players []Player, mode string) int {
	teamMap := make(map[int]int)
	for _, player := range players {
		teamMap[player.Team]++
	}
	teamMax := 0
	switch mode {
	case GameModePersonal:
		for i := 1; i < 100; i++ {
			if _, ok := teamMap[i]; !ok {
				return i
			}
		}
		return 0
	case GameModeTogether:
		return 1
	case GameModeTeamTwo:
		teamMax = 2
	case GameModeTeamThree:
		teamMax = 3
	case GameModeTeamFour:
		teamMax = 4
	}
	if teamMax == 0 {
		return 0
	}
	for i := 1; i < 100; i++ {
		if v, ok := teamMap[i]; ok {
			if v < teamMax {
				return i
			}
		} else {
			return i
		}
	}
	return 0
}

// GetPlayers 获取房间内玩家信息
func (s *roomService) GetPlayers(roomID int) (info []PlayerInfo, err error) {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return
	}
	for _, player := range room.Players {
		info = append(info, PlayerInfo{
			Player: player,
			Info:   s.Service.User.GetUserBaseInfo(player.UserID),
		})
	}
	return
}

// AddRoom 新建一个房间
func (s *roomService) AddRoom(ownID, title, mode, gameMap, password string, maxPlayer int) (roomID int, err error) {
	if maxPlayer > 20 {
		return 0, ErrMaxPlayer
	}
	roomID = s.newRoomID()
	s.Rooms = append(s.Rooms, GameRoom{
		ID:        roomID,
		Port:      -1,
		GameMap:   gameMap,
		MaxPlayer: maxPlayer,
		Password:  password,
		Title:     title,
		Mode:      mode,
		OwnID:     ownID,
		Playing:   false,
		Players: []Player{{
			UserID:  ownID,
			GameID:  0,
			IsReady: true,
			Team:    1,
			RoleID:  "new",
		}},
	})
	return
}

// JoinRoom 加入房间
func (s *roomService) JoinRoom(roomID int, userID, password string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	if room.Password != password {
		return ErrPassword
	}
	if room.MaxPlayer == len(room.Players) {
		return ErrMaxPlayer
	}
	room.Players = append(room.Players, Player{
		UserID:  userID,
		GameID:  s.newGameID(room.Players),
		RoleID:  "new",
		IsReady: false,
		Team:    s.getTeamID(room.Players, room.Mode),
	})
	return nil
}

// StartGame 开始游戏
func (s *roomService) StartGame(roomID int, ownID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	// 房主权限
	if room.OwnID != ownID {
		return ErrNotAllow
	}
	// 玩家是否全部已经准备
	for _, p := range room.Players {
		if p.IsReady == false {
			return ErrNotReady
		}
	}
	// 检测非合作模式是否全部都是一队的
	if room.Mode != GameModeTogether {
		team := room.Players[0].Team
		isOne := true
		for _, player := range room.Players {
			if player.Team != team {
				isOne = false
			}
		}
		if isOne {
			return ErrOneTeam
		}
	}

	// 建立房间服务器
	room.Port = s.Game.NewRoom(len(room.Players))
	if room.Port == -1 {
		// 服务器已满
		return ErrMaxServer
	}
	room.Playing = true
	return nil
}

// SetReady 设置准备状态
func (s *roomService) SetReady(roomID int, userID string, isReady bool) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	for i := range room.Players {
		if room.Players[i].UserID == userID {
			room.Players[i].IsReady = isReady
			return nil
		}
	}
	// 找不到用户
	return ErrNotFound
}

// SetReady 设置角色
func (s *roomService) SetRole(roomID int, userID, role string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	for i := range room.Players {
		if room.Players[i].UserID == userID {
			room.Players[i].RoleID = role
			return nil
		}
	}
	// 找不到用户
	return ErrNotFound
}

// SetTeam 设置队伍
func (s *roomService) SetTeam(roomID, teamID int, userID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	teamMax := 0
	switch room.Mode {
	case GameModePersonal:
		return ErrNotAllow
	case GameModeTogether:
		return ErrNotAllow
	case GameModeTeamTwo:
		teamMax = 2
	case GameModeTeamThree:
		teamMax = 3
	case GameModeTeamFour:
		teamMax = 4
	}
	if teamMax == 0 {
		return ErrNotAllow
	}

	teamMap := make(map[int]int)
	for _, player := range room.Players {
		teamMap[player.Team]++
	}
	if v, ok := teamMap[teamID]; ok && v >= teamMax {
		return ErrMaxPlayer
	}

	for i := range room.Players {
		if room.Players[i].UserID == userID {
			room.Players[i].Team = teamID
			break
		}
	}
	return nil
}

// GetOutRoom 房主踢人
func (s *roomService) GetOutRoom(roomID int, ownID, userID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	if room.OwnID != ownID {
		return ErrNotAllow
	}
	return s.QuitRoom(roomID, userID)
}

// SetRoomOwn 设置新的房主
func (s *roomService) SetRoomOwn(roomID int, ownID, newOwnID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	if room.OwnID != ownID {
		return ErrNotAllow
	}
	for _, player := range room.Players {
		if player.UserID == newOwnID {
			room.OwnID = newOwnID
			return nil
		}
	}
	return ErrNotFound
}

// SetRoomInfo 设置房间信息
func (s *roomService) SetRoomInfo(roomID, maxPlayer int, ownID, gameMap, gameMode string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	if room.OwnID != ownID {
		return ErrNotAllow
	}
	if gameMap != "" {
		room.GameMap = gameMap
	}
	if gameMode != "" {
		room.GameMap = gameMode
	}
	if maxPlayer != 0 {
		room.MaxPlayer = maxPlayer
	}
	return nil
}

func (s *roomService) deleteRoom(roomID int) error {
	for i := range s.Rooms {
		if s.Rooms[i].ID == roomID {
			s.Rooms = append(s.Rooms[:i], s.Rooms[i+1:]...)
		}
	}
	return ErrNotFound
}

// QuitRoom 退出房间
func (s *roomService) QuitRoom(roomID int, userID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	for i := range room.Players {
		if room.Players[i].UserID == userID {
			if len(room.Players) == 1 {
				// 最后一个人退出
				s.deleteRoom(roomID)
			} else {
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
				if userID == room.OwnID {
					// 如果房主走了传递房主权限
					room.OwnID = room.Players[0].UserID
				}
			}
			return nil
		}
	}
	return ErrNotFound
}
