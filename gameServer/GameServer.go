package gameServer

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// ServerConfig ...
type ServerConfig struct {
	PortPoolBegin int `yaml:"PortPoolBegin"` // 游戏服务器地址池开始
	PortPoolSize  int `yaml:"PortPoolSize"`  // 最大负载
}

// GameServer ...
type GameServer struct {
	Config ServerConfig  // 配置
	Room   []RoomData    // 房间
	Lock   *sync.RWMutex // 房间读写锁
}

// RoomData ...
type RoomData struct {
	Using        bool // 是否使用中
	Port         int
	Running      bool          // 是否已经开始
	Conn         *net.UDPConn  // 连接会话
	Players      []Player      // 房间玩家
	Frame        []FrameState  // 房间帧
	CurrentFrame int           // 当前帧
	MaxPeople    int           // 人数上限
	CreateTime   time.Time     // 创建时间
	Lock         *sync.RWMutex // 读写锁
}

// Player ...
type Player struct {
	Addr      *net.UDPAddr // 玩家地址
	ID        int          // 玩家ID
	Frame     int          // 玩家当前帧
	MissFrame int          // 玩家丢失帧
}

// InitServer 初始化游戏服务器
func (s *GameServer) InitServer(c ServerConfig) {
	s.Config = c
	s.Room = make([]RoomData, s.Config.PortPoolSize)
	s.Lock = new(sync.RWMutex)
}

func (s *GameServer) clearRoom() bool {
	for i := range s.Room {
		if s.Room[i].Using == false {
			continue
		}
		if (s.Room[i].Running == false || len(s.Room[i].Players) == 0) && time.Now().Unix()-s.Room[i].CreateTime.Unix() > 60 {
			s.closeRoom(i)
			return false
		}
	}
	return true
}

func (s *GameServer) IsUsing(port int) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for _, room := range s.Room {
		if room.Port == port {
			return true
		}
	}
	return false
}

func (s *GameServer) closeRoom(id int) {
	s.Room[id].Conn.Close()
	s.Room[id] = RoomData{}
	fmt.Println("close room ", id)
}

func (s *GameServer) getNullRoom() int {
	for i := range s.Room {
		if s.Room[i].Using == false {
			return i
		}
	}
	return -1
}

// NewRoom 开房
func (s *GameServer) NewRoom(maxPeople int) (port int) {
	s.Lock.Lock()
	s.clearRoom()
	roomID := s.getNullRoom()
	if roomID == -1 {
		s.Lock.Unlock()
		// 负载以达上限
		return -1
	}
	port = s.Config.PortPoolBegin + roomID
	service := ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Println("ResolveUDPAddr Error: ", err.Error())
		s.Lock.Unlock()
		return -1
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("ListenUDP Error: ", err.Error())
		s.Lock.Unlock()
		return -1
	}
	fmt.Println("GameServer is running in " + service)
	room := &s.Room[roomID]
	room.Conn = conn
	room.Port = port
	room.MaxPeople = maxPeople
	room.Lock = new(sync.RWMutex)
	room.CreateTime = time.Now()
	room.Using = true
	s.Lock.Unlock()
	go s.handleClient(roomID)
	return
}
