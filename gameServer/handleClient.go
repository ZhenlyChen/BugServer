package gameServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

// UserData ...
type UserData struct {
	ID    int     `json:"i"`
	Input int     `json:"c"`
	LocX  float32 `json:"x"`
	LocY  float32 `json:"y"`
	Dir   int     `json:"d"`
}

// UserBack ...
type UserBack struct {
	ID    int `json:"i"`
	Frame int `json:"f"`
}

// UserComeIn ...
type UserComeIn struct {
	ID  int    `json:"i"`
	Key string `json:"k"`
}

func (s *GameServer) handleClient(id int) {
	for {
		var buf [1024]byte
		_, addr, err := s.Room[id].Conn.ReadFromUDP(buf[0:]) // 等待连接
		if buf[1023] != 0 {
			fmt.Println("Too Long Data")
			continue
		}
		if err != nil || s.Room[id].Using == false {
			fmt.Println("Error data from room ", id)
			return
		}
		if buf[0] == '0' { // 加入房间
			s.joinRoom(id, &buf, addr)
		} else if buf[0] == '1' { // 传入数据
			s.setInput(id, &buf)
		} else if buf[0] == '2' { // 设置帧数
			s.setFrame(id, &buf)
		} else if buf[0] == '3' { // 退出房间
			s.goOutRoom(id, addr)
			// 删除对局
			if len(s.Room[id].Players) == 0 {
				s.Lock.Lock()
				s.closeRoom(id)
				s.Lock.Unlock()
			}
			break
		}
	}
}

func (s *GameServer) joinRoom(id int, buf *[1024]byte, addr *net.UDPAddr) {
	data := UserComeIn{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		s.Room[id].Lock.Lock()
		room := &s.Room[id]
		if room.Using == false {
			// 房间已关闭
			room.Lock.Unlock()
			return
		}
		if len(room.Players) >= room.MaxPeople {
			// 人数已满
			room.Lock.Unlock()
			return
		}
		for i := range room.Players {
			if room.Players[i].ID == data.ID || room.Players[i].Addr.String() == addr.String() {
				// 已经加入
				room.Conn.WriteToUDP(append([]byte("join"), 0), addr)
				room.Lock.Unlock()
				return
			}
		}
		room.Players = append(room.Players, Player{
			Addr:      addr,
			ID:        data.ID,
			Frame:     0,
			MissFrame: 0,
		})
		room.Lock.Unlock()
		fmt.Println("Come in ", addr.String())
		fmt.Println(len(room.Players), "/", room.MaxPeople)
		room.Conn.WriteToUDP(append([]byte("join"), 0), addr)
		if len(room.Players) == room.MaxPeople && room.Running == false {
			fmt.Println("Game Begin")
			// 开始发送帧信息
			go s.sendAll(id)
		}
	}
}

func (s *GameServer) setInput(id int, buf *[1024]byte) {
	if !s.Room[id].Running || !s.Room[id].Using {
		return
	}
	data := UserData{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		// 写入帧，互斥锁
		s.Room[id].Lock.Lock()
		room := &s.Room[id]
		currentFrame := room.CurrentFrame - 1
		room.Frame[currentFrame].Commands = append(room.Frame[currentFrame].Commands, Command{
			UserID: data.ID,
			Input:  data.Input,
			LocX:   data.LocX,
			LocY:   data.LocY,
			Dir:    data.Dir,
		})
		room.Lock.Unlock()
	}
}

func (s *GameServer) setFrame(id int, buf *[1024]byte) {
	data := UserBack{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		room := &s.Room[id]
		for i := range room.Players {
			if room.Players[i].ID == data.ID {
				room.Lock.Lock()
				room.Players[i].Frame = data.Frame
				room.Players[i].MissFrame = 0
				room.Lock.Unlock()
				break
			}
		}
	}
}

func (s *GameServer) goOutRoom(id int, addr *net.UDPAddr) {
	s.Room[id].Lock.Lock()
	room := &s.Room[id]
	for i := range room.Players {
		if room.Players[i].Addr.String() == addr.String() {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			fmt.Println("Go out: ", addr.String())
			room.Conn.WriteToUDP(append([]byte("out"), 0), addr)
			break
		}
	}
	room.Lock.Unlock()
}
