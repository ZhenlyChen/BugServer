package gameServer

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResData ...
type ResData struct {
	Data []FrameState `json:"d"`
}

// FrameState ...
type FrameState struct {
	FrameID  int       `json:"f"`
	Commands []Command `json:"c"`
}

// Commend ...
type Command struct {
	UserID int     `json:"i"`
	Input  int     `json:"c"`
	LocX   float32 `json:"x"`
	LocY   float32 `json:"y"`
	Dir    int     `json:"d"`
}

// 并发发送数据
func (s *GameServer) sendToPlayer(rID, pID int, c chan int) {
	var res ResData
	// 检测是否掉线
	room := &s.Room[rID]
	player := &room.Players[pID]
	if player.MissFrame > 200 {
		// 判断为已经掉线
		player.MissFrame = 999
		c <- 0
		return
	}
	player.MissFrame++
	if room.CurrentFrame-player.Frame > 5 {
		res = ResData{
			Data: room.Frame[player.Frame : player.Frame+5],
		}
	} else {
		res = ResData{
			Data: room.Frame[player.Frame:],
		}
	}
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	room.Conn.WriteToUDP(b, player.Addr)
	c <- 0
}

func (s *GameServer) sendAll(id int) {
	for {
		time.Sleep(time.Millisecond * 67)
		s.Lock.Lock()
		room := &s.Room[id]
		if room.Using == false {
			s.Lock.Unlock()
			return
		}
		for i := range room.Players {
			if room.Players[i].MissFrame == 999 {
				// 清理掉线用户
				s.goOutRoom(id, room.Players[i].Addr)
				break
			}
		}
		if len(room.Frame) > 0 {

			// 并发发送数据给用户
			c := make(chan int)
			playerCount := len(room.Players)

			room.Lock.RLock()
			for i := 0; i < playerCount; i++ {
				go s.sendToPlayer(id, i, c)
			}
			for i := 0; i < playerCount; i++ {
				<-c
			}
			room.Lock.RUnlock()
			close(c)
		}
		// 增加新的帧， 互斥锁
		room.Lock.Lock()
		room.Frame = append(room.Frame, FrameState{
			FrameID:  room.CurrentFrame + 1,
			Commands: []Command{},
		})
		room.CurrentFrame++
		room.Running = true
		room.Lock.Unlock()
		s.Lock.Unlock()
		// 解锁
	}
}
