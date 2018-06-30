package services

import (
	"github.com/ZhenlyChen/BugServer/gameServer"
	"fmt"
)

// RoomService ...
type RoomService interface {
	InitGameServer(config gameServer.ServerConfig)
	NewRoom()
}

type roomService struct {
	Service  *Service
	Game *gameServer.GameServer
}

func (s *roomService) InitGameServer(config gameServer.ServerConfig) {
	s.Game = new(gameServer.GameServer)
	s.Game.InitServer(config)
}

func (s *roomService) NewRoom() {
	fmt.Println("New a room at ", s.Game.NewRoom(1))
}