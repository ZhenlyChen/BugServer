package gameServer

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

// ServerConfig ...
type ServerConfig struct {
	PortPoolBegin int `yaml:"PortPoolBegin"` // 游戏服务器地址池开始
	PortPoolSize  int `yaml:"PortPoolSize"`  // 最大负载
}

// GameServer ...
type GameServer struct {
	Config      ServerConfig
	CurrentLoad int
	Room        []RoomData
}

// RoomData ...
type RoomData struct {
	conn         *net.UDPConn
	Players      []Player
	Frame        []FrameState
	CurrentFrame int
	People       int
	Running      bool
	Lock         *sync.RWMutex
}
// Player ...
type Player struct {
	IP        *net.UDPAddr
	ID        int
	Frame     int
	MissFrame int
}

// InitServer 初始化游戏服务器
func (s *GameServer) InitServer(c ServerConfig) {
	s.Config = c
	s.CurrentLoad = 0
}

// NewRoom 开房
func (s *GameServer) NewRoom(people int) (port int) {
	if s.CurrentLoad > s.Config.PortPoolSize {
		// 负载以达上限
		return -1
	}
	service := ":" + strconv.Itoa(s.Config.PortPoolBegin+s.CurrentLoad)
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	fmt.Println("GameServer is running in " + service)
	s.Room = append(s.Room, RoomData{
		Players:      []Player{},
		Frame:        []FrameState{},
		CurrentFrame: 0,
		People:       people,
		Running:      false,
		Lock:         new(sync.RWMutex),
	})
	fmt.Println(s.Room)
	go s.handleClient(conn, s.CurrentLoad)
	s.CurrentLoad++
	return s.Config.PortPoolBegin + s.CurrentLoad - 1
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
