package services

import (
	"time"

	"github.com/ZhenlyChen/BugServer/gameServer"
	"math/rand"
	"sync"
)

// RoomService ...
type RoomService interface {
	InitGameServer(config gameServer.ServerConfig)
	IsInRoom(userID string) bool
	CheckHeart()
	Heart(userID string, roomID int) bool
	GetRoom(roomID int) (room *Room, err error)
	JoinRoom(roomID int, userID, password string) error
	SetReady(roomID int, userID string, isReady bool) error
	SetRole(roomID, roleID int, userID string) error
	SetTeam(roomID, teamID int, userID string) error
	GetRooms() []Room
	QuitRoom(room *Room, userID string) error
	// 房主
	StartGame(roomID int, ownID string) error
	SetRoomOwn(roomID int, ownID, newOwnID string) error
	SetPlaying(roomID int, userID string, isPlaying bool) error
	AddRoom(ownID, title, mode, gameMap, password string, maxPlayer int, isRandom bool) (roomID int, err error)
	SetRoomInfo(roomID, maxPlayer int, ownID, gameMap, title, password string, isRandom bool) error
	GetOutRoom(roomID int, ownID, userID string) error
}

type roomService struct {
	Service *Service
	Game    *gameServer.GameServer
	Rooms   []Room
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
	MaxRole   = 4
)

type Room struct {
	Using bool          // 是否使用中
	Lock  *sync.RWMutex // 读写锁
	Info  RoomInfo      // 信息
}

// RoomInfo 房间数据
type RoomInfo struct {
	ID          int          `json:"id"`          // 房间ID
	OwnID       string       `json:"ownId"`       // 房主ID
	OwnInfo     UserBaseInfo `json:"ownInfo"`     // 房主信息
	Port        int          `json:"port"`        // 房间服务器端口
	Title       string       `json:"title"`       // 标题
	IsRandom    bool         `json:"isRandom"`    // 是否随机角色
	GameMap     string       `json:"gameMap"`     // 游戏地图
	RandSeed    int          `json:"randSeed"`    // 随机种子
	MaxPlayer   int          `json:"maxPlayer"`   // 最大人数
	PlayerCount int          `json:"playerCount"` // 当前玩家数(传输时设置)
	Mode        string       `json:"mode"`        // 游戏模式
	Password    string       `json:"password"`    // 房间密码
	Playing     bool         `json:"playing"`     // 是否正在玩
	Players     []Player     `json:"players"`     // 玩家数据
}

// Player 玩家信息
type Player struct {
	UserID  string       `json:"userId"`  // 玩家ID
	Info    UserBaseInfo `json:"info"`    // 玩家信息
	GameID  int          `json:"gameId"`  // 游戏内ID
	RoleID  int          `json:"roleId"`  // 角色ID
	IsReady bool         `json:"isReady"` // 是否准备
	Heart   int          `json:"heart"`   // 心跳💗
	Team    int          `json:"team"`    // 玩家队伍
}

func (s *roomService) InitGameServer(config gameServer.ServerConfig) {
	s.Game = new(gameServer.GameServer)
	s.Game.InitServer(config)
}

func (s *roomService) CheckHeart() {
	for {
		time.Sleep(time.Second)
		for i := range s.Rooms {
			room := &s.Rooms[i]
			room.Lock.Lock()
			if room.Using == false || room.Info.Playing == true {
				continue
			}
			for j := range room.Info.Players {
				if room.Info.Players[j].Heart > 3 {
					room.Lock.Unlock()
					s.QuitRoom(room, room.Info.Players[j].UserID)
					room.Lock.Lock()
					break
				} else {
					room.Info.Players[j].Heart++
				}
			}
			room.Lock.Unlock()
		}
	}
}

func (s *roomService) Heart(userID string, roomID int) bool {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return false
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	for i := range room.Info.Players {
		if room.Info.Players[i].UserID == userID {
			room.Info.Players[i].Heart = 0
			return true
		}
	}
	return false
}

// GetRooms 获取房间列表
func (s *roomService) GetRooms() (rooms []Room) {
	for _, room := range s.Rooms {
		if room.Using == true {
			room.Info.Playing = s.Game.IsUsing(room.Info.Port)
			rooms = append(rooms, room)
		}
	}
	return
}

// findRoom 寻找房间，
func (s *roomService) GetRoom(roomID int) (room *Room, err error) {
	if roomID >= len(s.Rooms) {
		return nil, ErrNotFound
	}
	room = &s.Rooms[roomID]
	if room.Using {
		room.Info.Playing = s.Game.IsUsing(room.Info.Port)
		return room, nil
	}
	err = ErrNotFound
	return
}

// getGameID 获取新的游戏ID
func (room *Room) newGameID() int {
	for i := 0; i < MaxPlayer; i++ {
		hasExist := false
		for _, player := range room.Info.Players {
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
	for i := range s.Rooms {
		if s.Rooms[i].Using == false {
			return i
		}
	}
	if len(s.Rooms) < MaxRoom {
		newID := len(s.Rooms)
		s.Rooms = append(s.Rooms, Room{
			Using: false,
			Lock:  new(sync.RWMutex),
			Info: RoomInfo{
				ID: newID,
			},
		})
		return newID
	}
	return -1
}

func (room *Room) newTeamID() int {
	room.Lock.RLock()
	defer room.Lock.RUnlock()
	teamMap := make(map[int]int)
	for _, player := range room.Info.Players {
		teamMap[player.Team]++
	}
	teamMax := 0
	switch room.Info.Mode {
	case GameModePersonal:
		return -1
	case GameModeTogether:
		return 1
	case GameModeTeamTwo:
		teamMax = 2
	case GameModeTeamThree:
		teamMax = 3
	case GameModeTeamFour:
		teamMax = 4
	}
	for i := 1; i < MaxPlayer; i++ {
		if v, ok := teamMap[i]; ok {
			if v < teamMax {
				return i
			}
		} else {
			return i
		}
	}
	return -1
}

func (s *roomService) IsInRoom(userID string) bool {
	for _, room := range s.GetRooms() {
		for _, player := range room.Info.Players {
			if player.UserID == userID {
				return true
			}
		}
	}
	return false
}

// AddRoom 新建一个房间
func (s *roomService) AddRoom(ownID, title, mode, gameMap, password string, maxPlayer int, isRandom bool) (roomID int, err error) {
	if maxPlayer > 20 {
		return 0, ErrMaxPlayer
	}
	// 玩家是否已经在房间内
	if s.IsInRoom(ownID) {
		return 0, ErrNotAllow
	}
	ownInfo := s.Service.User.GetUserBaseInfo(ownID)
	roomID = s.newRoomID()
	newRoom := &s.Rooms[roomID]
	newRoom.Using = true
	newRoom.Info = RoomInfo{
		ID:        roomID,
		Port:      -1,
		GameMap:   gameMap,
		MaxPlayer: maxPlayer,
		Password:  password,
		Title:     title,
		Mode:      mode,
		OwnID:     ownID,
		OwnInfo:   ownInfo,
		Playing:   false,
		IsRandom:  isRandom,
		Players: []Player{{
			UserID:  ownID,
			Info:    ownInfo,
			GameID:  0,
			IsReady: true,
			Team:    1,
			RoleID:  0,
		}},
	}
	return
}

// JoinRoom 加入房间
func (s *roomService) JoinRoom(roomID int, userID, password string) error {
	if s.IsInRoom(userID) {
		return ErrNotAllow
	}
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.Info.Password != password {
		return ErrPassword
	}
	if room.Info.MaxPlayer >= len(room.Info.Players) {
		return ErrMaxPlayer
	}
	room.Info.Players = append(room.Info.Players, Player{
		UserID:  userID,
		Info:    s.Service.User.GetUserBaseInfo(userID),
		GameID:  room.newGameID(),
		RoleID:  0,
		IsReady: false,
		Team:    room.newTeamID(),
	})
	return nil
}

// StartGame 开始游戏
func (s *roomService) StartGame(roomID int, ownID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	// 房主权限
	if room.Info.OwnID != ownID {
		return ErrNotAllow
	}
	// 是否在玩
	if room.Info.Playing || room.Info.Port != -1 {
		return ErrNotAllow
	}
	// 若开始失败重置开始状态，各个用户返回房间列表
	room.Info.Playing = false
	// 玩家是否全部已经准备
	for _, p := range room.Info.Players {
		if p.IsReady == false {
			return ErrNotReady
		}
	}
	// 检测非合作模式是否全部都是一队的
	if room.Info.Mode != GameModeTogether && room.Info.Mode != GameModePersonal {
		teamMap := make(map[int]int)
		for _, player := range room.Info.Players {
			teamMap[player.Team]++
		}
		if len(teamMap) <= 1 {
			return ErrOneTeam
		}
	}
	// 随机分配角色
	if room.Info.IsRandom {
		gameInfo := s.Service.Game.GetNewestVersion()
		for i := range room.Info.Players {
			room.Info.Players[i].RoleID = rand.Intn(gameInfo.MaxRole)
		}
	}
	// 生成随机数种子
	room.Info.RandSeed = rand.Intn(167167167)
	// 建立房间服务器
	room.Info.Port = s.Game.NewRoom(len(room.Info.Players))
	if room.Info.Port == -1 {
		// 服务器已满
		return ErrMaxServer
	}
	room.Info.Playing = true
	return nil
}

// SetReady 设置准备状态
func (s *roomService) SetReady(roomID int, userID string, isReady bool) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if userID == room.Info.OwnID {
		// 房主不能改变准备状态
		return ErrNotAllow
	}
	for i := range room.Info.Players {
		if room.Info.Players[i].UserID == userID {
			room.Info.Players[i].IsReady = isReady
			return nil
		}
	}
	// 找不到用户
	return ErrNotFound
}

// SetPlaying 设置开始状态
func (s *roomService) SetPlaying(roomID int, userID string, isPlaying bool) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.Info.OwnID != userID {
		return ErrNotAllow
	}
	room.Info.Playing = isPlaying
	// 找不到用户
	return nil
}

// SetReady 设置角色
func (s *roomService) SetRole(roomID, roleID int, userID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	for i := range room.Info.Players {
		if room.Info.Players[i].UserID == userID {
			room.Info.Players[i].RoleID = roleID
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
	room.Lock.Lock()
	defer room.Lock.Unlock()
	teamMax := 0
	switch room.Info.Mode {
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
	for _, player := range room.Info.Players {
		teamMap[player.Team]++
	}
	if v, ok := teamMap[teamID]; ok && v >= teamMax {
		return ErrMaxPlayer
	}

	for i := range room.Info.Players {
		if room.Info.Players[i].UserID == userID {
			room.Info.Players[i].Team = teamID
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
	room.Lock.RLock()
	if room.Info.OwnID != ownID {
		return ErrNotAllow
	}
	room.Lock.RUnlock()
	return s.QuitRoom(room, userID)
}

// SetRoomOwn 设置新的房主
func (s *roomService) SetRoomOwn(roomID int, ownID, newOwnID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.Info.OwnID != ownID {
		return ErrNotAllow
	}
	for _, player := range room.Info.Players {
		if player.UserID == newOwnID {
			room.Info.OwnID = newOwnID
			room.Info.OwnInfo = s.Service.User.GetUserBaseInfo(newOwnID)
			return nil
		}
	}
	return ErrNotFound
}

// SetRoomInfo 设置房间信息
func (s *roomService) SetRoomInfo(roomID, maxPlayer int, ownID, gameMap, title, password string, isRandom bool) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.Info.OwnID != ownID {
		return ErrNotAllow
	}
	if len(room.Info.Players) > maxPlayer {
		return ErrNotAllow
	}
	room.Info.GameMap = gameMap
	room.Info.MaxPlayer = maxPlayer
	room.Info.IsRandom = isRandom
	room.Info.Password = password
	room.Info.Title = title
	return nil
}

// QuitRoom 退出房间
func (s *roomService) QuitRoom(room *Room, userID string) error {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	for i := len(room.Info.Players) - 1; i >= 0; i-- {
		if room.Info.Players[i].UserID == userID {
			if len(room.Info.Players) == 1 {
				// 最后一个人退出
				room.Using = false
				room.Info = RoomInfo{}
			} else {
				room.Info.Players = append(room.Info.Players[:i], room.Info.Players[i+1:]...)
				if userID == room.Info.OwnID {
					// 如果房主走了传递房主权限
					room.Info.OwnID = room.Info.Players[0].UserID
				}
			}
			return nil
		}
	}
	return ErrNotFound
}
