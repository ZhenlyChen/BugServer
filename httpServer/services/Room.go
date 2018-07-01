package services

import (
	"github.com/ZhenlyChen/BugServer/gameServer"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
)

// RoomService ...
type RoomService interface {
	InitGameServer(config gameServer.ServerConfig)
	JoinRoom(roomID, userID string) error
	SetReady(roomID, userID string, isReady bool) error
	SetRole(roomID, userID string, role string) error
	GetRooms() []gameRoom
	GetPlayers(roomID string) (info []PlayerInfo, err error)
	// 房主
	StartGame(roomID, ownID string) error
	AddRoom(ownID, mode, gameMap string, maxPlayer int) (roomID string)
	SetRoomInfo(roomID, gameMap, gameMode string, maxPlayer int) error
}

type roomService struct {
	Service *Service
	Game    *gameServer.GameServer
	Rooms   []gameRoom
}

const (
	GamemodePersonal  = "personal" // 个人
	GamemodeTogether  = "together" // 合作
	GamemodeTeamtwo   = "team2"    // 2人团队
	GamemodeTeamthree = "team3"    // 3人团队
	GamemodeTeamfour  = "team4"    // 4人团队
)

var (
	ErrNotfound  = errors.New("not_found")
	ErrNotallow  = errors.New("not_allow")
	ErrNotready  = errors.New("not_ready")
	ErrServermax = errors.New("server_max")
	ErrMaxplayer = errors.New("max_player")
)

type gameRoom struct {
	ID        string
	OwnID     string // 房主ID
	Port      int
	GameMap   string
	MaxPlayer int
	Mode      string
	Playing   bool
	Players   []player
}

type player struct {
	UserID  string
	GameID  int
	RoleID  string // 角色ID
	IsReady bool
	Team    int // "1-4" - 队伍一~四
}

func (s *roomService) InitGameServer(config gameServer.ServerConfig) {
	s.Game = new(gameServer.GameServer)
	s.Game.InitServer(config)
}

// GetRooms 获取房间列表
func (s *roomService) GetRooms() []gameRoom {
	return s.Rooms
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	Player player
	Info   UserBaseInfo
}

// findRoom 修改房间，返回房间地址
func (s *roomService) findRoom(roomID string) (room *gameRoom, err error) {
	for i := range s.Rooms {
		if roomID == s.Rooms[i].ID {
			room = &s.Rooms[i]
			return
		}
	}
	err = ErrNotfound
	return
}

// getGameID 获取新的游戏ID
func (s *roomService) newGameID(players []player) int {
	for i := 0; i < 100; i++ {
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
	return 100
}

func (s *roomService) getTeamID(players []player, mode string) int {
	teamMap := make(map[int]int)
	for _, player := range players {
		teamMap[player.Team]++
	}
	teamMax := 0
	switch mode {
	case GamemodePersonal:
		for i := 1; i < 100; i++ {
			if _, ok := teamMap[i]; !ok {
				return i
			}
		}
		return 0
	case GamemodeTogether:
		return 1
	case GamemodeTeamtwo:
		teamMax = 2
	case GamemodeTeamthree:
		teamMax = 3
	case GamemodeTeamfour:
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
func (s *roomService) GetPlayers(roomID string) (info []PlayerInfo, err error) {
	room, err := s.findRoom(roomID)
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
func (s *roomService) AddRoom(ownID, mode, gameMap string, maxPlayer int) (roomID string, err error) {
	if maxPlayer > 20 {
		return "", ErrMaxplayer
	}
	roomID = bson.NewObjectId().Hex()
	s.Rooms = append(s.Rooms, gameRoom{
		ID:        roomID,
		Port:      -1,
		GameMap:   gameMap,
		MaxPlayer: maxPlayer,
		Mode:      mode,
		OwnID:     ownID,
		Playing:   false,
		Players: []player{{
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
func (s *roomService) JoinRoom(roomID, userID string) error {
	room, err := s.findRoom(roomID)
	if err != nil {
		return err
	}
	if room.MaxPlayer == len(room.Players) {
		return ErrMaxplayer
	}
	room.Players = append(room.Players, player{
		UserID:  userID,
		GameID:  s.newGameID(room.Players),
		RoleID:  "new",
		IsReady: false,
		Team:    s.getTeamID(room.Players, room.Mode),
	})
	return nil
}

// StartGame 开始游戏
func (s *roomService) StartGame(roomID, ownID string) error {
	room, err := s.findRoom(roomID)
	if err != nil {
		return err
	}
	// 房主权限
	if room.OwnID != ownID {
		return ErrNotallow
	}
	// 玩家是否全部已经准备
	for _, p := range room.Players {
		if p.IsReady == false {
			return ErrNotready
		}
	}
	// 建立房间服务器
	room.Port = s.Game.NewRoom(len(room.Players))
	if room.Port == -1 {
		// 服务器已满
		return ErrServermax
	}
	room.Playing = true
	return nil
}

// SetReady 设置准备状态
func (s *roomService) SetReady(roomID, userID string, isReady bool) error {
	room, err := s.findRoom(roomID)
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
	return ErrNotfound
}

// SetReady 设置角色
func (s *roomService) SetRole(roomID, userID string, role string) error {
	room, err := s.findRoom(roomID)
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
	return ErrNotfound
}


// SetRoomInfo 设置房间信息
func (s *roomService) SetRoomInfo(roomID, gameMap, gameMode string, maxPlayer int) error {
	room, err := s.findRoom(roomID)
	if err != nil {
		return err
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
}