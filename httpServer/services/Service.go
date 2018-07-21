package services

import (
	"errors"

	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"sync"
)

// Service ...
type Service struct {
	Model *models.Model
	User  userService
	Room  roomService
	Game  gameService
}

// Err 错误类型
var (
	ErrNotFound  = errors.New("not_found")
	ErrNotAllow  = errors.New("not_allow")
	ErrNotReady  = errors.New("not_ready")
	ErrMaxServer = errors.New("max_server")
	ErrMaxPlayer = errors.New("max_player")
	ErrPassword  = errors.New("err_password")
	ErrOneTeam   = errors.New("one_team")
)

// NewService ...
func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	service.User = userService{
		Model:    &m.User,
		Service:  service,
		UserInfo: make(map[string]UserBaseInfo),
	}
	service.Room = roomService{
		Service: service,
		roomLock: new(sync.RWMutex),
	}
	service.Game = gameService{
		Service: service,
		Model:   &m.Game,
	}
	return service
}

// GetUserService ...
func (s *Service) GetUserService() UserService {
	return &s.User
}

// GetRoomService ...
func (s *Service) GetRoomService() RoomService {
	return &s.Room
}

// GetGameService ...
func (s *Service) GetGameService() GameService {
	return &s.Game
}
