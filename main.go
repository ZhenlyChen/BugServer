package main

import (
	"net/http"

	"github.com/ZhenlyChen/BugServer/router"
	_ "github.com/ZhenlyChen/BugServer/router"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/peer"
	httppeer "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
)

func main() {
	queue := cellnet.NewEventQueue()


	p := peer.NewGenericPeer("http.Acceptor", "server", "127.0.0.1:18801", queue)

	proc.BindProcessorHandler(p, "http", func(raw cellnet.Event) {
		switch msg := raw.Message().(type) {
		case *router.HttpLoginREQ:
			println(msg.UserName)
			println(msg.Password)
			raw.Session().Send(&httppeer.MessageRespond{
				StatusCode: http.StatusOK,
				Msg: &router.HttpTokenACK{
					Token: "ok",
				},
			})
		}
	})

	p.Start()

	queue.StartLoop()

	queue.Wait()
}
