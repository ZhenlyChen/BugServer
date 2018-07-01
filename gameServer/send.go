package gameServer

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResData ...
type ResData struct {
	Data []FrameState `json:"data"`
}

// FrameState ...
type FrameState struct {
	FrameID  int       `json:"frameID"`
	Commends []Commend `json:"commends"`
}

// Commend ...
type Commend struct {
	UserID int `json:"id"`
	Input  int `json:"input"`
}

// 并发发送数据
func (s *GameServer) sendToPlayer(rID, pID int, c chan int) {
	var res ResData
	// 检测是否掉线
	fmt.Println(s.Room[rID].Players[pID].MissFrame)
	if s.Room[rID].Players[pID].MissFrame > 100 {
		// 判断为已经掉线
		s.goOutRoom(rID, s.Room[rID].Players[pID].IP)
		c <- 0
		return
	}
	s.Room[rID].Players[pID].MissFrame++
	// 读取帧数据, 共享锁
	s.Room[rID].Lock.RLock()
	if s.Room[rID].CurrentFrame-s.Room[rID].Players[pID].Frame > 10 {
		res = ResData{
			Data: s.Room[rID].Frame[s.Room[rID].Players[pID].Frame : s.Room[rID].Players[pID].Frame+10],
		}
	} else {
		res = ResData{
			Data: s.Room[rID].Frame[s.Room[rID].Players[pID].Frame:],
		}
	}
	s.Room[rID].Lock.RUnlock()
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	s.Room[rID].conn.WriteToUDP(b, s.Room[rID].Players[pID].IP)
	c <- 0
}

func (s *GameServer) sendAll(id int) {
	s.Room[id].Running = true
	for {
		// 并发发送数据给用户
		c := make(chan int)
		playerCount := len(s.Room[id].Players)
		for i := 0; i < playerCount; i++ {
			go s.sendToPlayer(id, i, c)
		}
		for i := 0; i < playerCount; i++ {
			<-c
		}
		close(c)
		// 增加新的帧， 互斥锁
		s.Room[id].Lock.Lock()
		s.Room[id].Frame = append(s.Room[id].Frame, FrameState{
			FrameID:  s.Room[id].CurrentFrame + 1,
			Commends: []Commend{},
		})
		s.Room[id].CurrentFrame++
		s.Room[id].Lock.Unlock()
		// 解锁
		time.Sleep(time.Millisecond * 100)
	}
}
