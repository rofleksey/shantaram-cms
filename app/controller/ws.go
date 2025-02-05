package controller

import (
	"github.com/gofiber/contrib/websocket"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
)

type WebSocket struct {
	wsService *service.WebSocket
}

func NewWebSocket(
	wsService *service.WebSocket,
) *WebSocket {
	return &WebSocket{
		wsService: wsService,
	}
}

func (c *WebSocket) GlobalHandler(conn *websocket.Conn) {
	done := conn.Locals("done").(<-chan struct{})

	sub := c.wsService.Subscribe(dao.WsChannelGlobal, func(data dao.WsMessage) {
		_ = conn.WriteJSON(data)
	})

	defer c.wsService.Unsubscribe(sub)

	<-done
}
