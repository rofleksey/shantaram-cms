package service

import (
	"github.com/simonfxr/pubsub"
	"shantaram-cms/app/dao"
)

type WebSocket struct {
	bus *pubsub.Bus
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		bus: pubsub.NewBus(),
	}
}

func (s *WebSocket) Subscribe(channel string, callback func(message dao.WsMessage)) *pubsub.Subscription {
	return s.bus.Subscribe(channel, callback)
}

func (s *WebSocket) Unsubscribe(sub *pubsub.Subscription) {
	s.bus.Unsubscribe(sub)
}

func (s *WebSocket) Publish(channel string, data dao.WsMessage) {
	s.bus.Publish(channel, data)
}
