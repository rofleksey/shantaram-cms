package dao

const WsChannelGlobal = "global"

const WsEventOrdersChanged = "orders_changed"

type WsMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
