package web

import (
	"bobthebuilder/builder"
	"github.com/cloudfoundry/gosigar"
	"golang.org/x/net/websocket"
	//"github.com/hoisie/web"
	"io"
	"time"
)

// Echo the data received on the WebSocket.
func ws_EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func ws_EventServer(ws *websocket.Conn) {
	ok := true

	eventMessages := make(chan builder.BuilderEvent, 10)
	builder.GetInstance().SubscribeToEvents(eventMessages)
	defer builder.GetInstance().UnsubscribeFromEvents(eventMessages)
	defer func() { ok = false }() //signal other routines to die

	go sendServerStatsMessages(eventMessages, &ok)
	for msg := range eventMessages {
		err := websocket.JSON.Send(ws, msg)
		if err != nil {
			return
		}
	}
}

func sendServerStatsMessages(msgBuffer chan builder.BuilderEvent, ok *bool) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		if !(*ok) {
			return
		}

		mem := sigar.Mem{}
		swap := sigar.Swap{}
		mem.Get()
		swap.Get()

		msgBuffer <- builder.BuilderEvent{
			Type: builder.EVT_SERVER_STATS,
			Data: map[string]interface{}{
				"mem":  mem,
				"swap": swap,
			},
			Index: -1,
		}
	}
}
