package web

import (
  "golang.org/x/net/websocket"
  "bobthebuilder/builder"
  //"github.com/hoisie/web"
  "io"
)


// Echo the data received on the WebSocket.
func ws_EchoServer(ws *websocket.Conn) {
    io.Copy(ws, ws)
}

func ws_EventServer(ws *websocket.Conn) {
  eventMessages := make(chan builder.BuilderEvent, 10)
  builder.GetInstance().SubscribeToEvents(eventMessages)
  defer builder.GetInstance().UnsubscribeFromEvents(eventMessages)

  for msg := range eventMessages {
    err := websocket.JSON.Send(ws, msg)
    if err != nil{
      return
    }
  }
}
