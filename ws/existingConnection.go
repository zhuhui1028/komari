package ws

import "github.com/gorilla/websocket"

var (
	ConnectedClients = make(map[string]*websocket.Conn)
	ConnectedUsers   = []*websocket.Conn{}
	LatestReport     = []map[string]interface{}{}
)
