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
	ID int `json:"i"`
}

func (s *GameServer) handleClient(conn *net.UDPConn, id int) {
	s.Room[id].conn = conn
	for {
		var buf [1024]byte
		_, addr, err := conn.ReadFromUDP(buf[0:]) // 等待连接
		if err != nil || buf[1023] != 0 || len(s.Room) <= id {
			fmt.Println("Error")
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
				s.Room = append(s.Room[:id], s.Room[id+1:]...)
			}
			break
		}
	}
}

func (s *GameServer) joinRoom(id int, buf *[1024]byte, addr *net.UDPAddr) {
	data := UserComeIn{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		for i := range s.Room[id].Players {
			if s.Room[id].Players[i].ID == data.ID || s.Room[id].Players[i].IP.String() == addr.String() {
				s.Room[id].conn.WriteToUDP([]byte("ok"), addr)
				return
			}
		}
		s.Room[id].Players = append(s.Room[id].Players, Player{
			IP:        addr,
			ID:        data.ID,
			Frame:     0,
			MissFrame: 0,
		})
		if len(s.Room[id].Players) == s.Room[id].People && s.Room[id].Running == false {
			fmt.Println("Game Begin")
			go s.sendAll(id)
		}
		fmt.Println("Come in ", addr.String())
		fmt.Println(len(s.Room[id].Players), '/', s.Room[id].People)
		s.Room[id].conn.WriteToUDP([]byte("join"), addr)
	}
}

func (s *GameServer) setInput(id int, buf *[1024]byte) {
	data := UserData{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		// 写入帧，互斥锁
		s.Room[id].Lock.Lock()
		currentFrame := s.Room[id].CurrentFrame - 1
		s.Room[id].Frame[currentFrame].Commands = append(s.Room[id].Frame[currentFrame].Commands, Command{
			UserID: data.ID,
			Input:  data.Input,
			LocX:   data.LocX,
			LocY:   data.LocY,
			Dir:    data.Dir,
		})
		s.Room[id].Lock.Unlock()
	}
}

func (s *GameServer) setFrame(id int, buf *[1024]byte) {
	data := UserBack{}
	if err := json.Unmarshal(buf[1:bytes.IndexByte(buf[1:], 0)+1], &data); err == nil {
		for i := range s.Room[id].Players {
			if s.Room[id].Players[i].ID == data.ID {
				s.Room[id].Players[i].Frame = data.Frame
				s.Room[id].Players[i].MissFrame = 0
				continue
			}
		}
	}
}

func (s *GameServer) goOutRoom(id int, addr *net.UDPAddr) {
	for i := range s.Room[id].Players {
		if s.Room[id].Players[i].IP.String() == addr.String() {
			s.Room[id].Players = append(s.Room[id].Players[:i], s.Room[id].Players[i+1:]...)
			fmt.Println("Go out: ", addr.String())
			s.Room[id].conn.WriteToUDP([]byte("out"), addr)
			break
		}
	}
}
