package controller

import (
	"bytes"
	"context"
	"log/slog"
	"shantaram/app/api"
	"shantaram/app/service/auth"
	"shantaram/app/service/pubsub"
	"shantaram/pkg/config"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/jellydator/ttlcache/v3"
	"github.com/samber/do"
)

var pingMsg = []byte("ping")

type WS struct {
	cfg           *config.Config
	authService   *auth.Service
	pubSubService *pubsub.Service
}

func NewWS(di *do.Injector) *WS {
	return &WS{
		cfg:           do.MustInvoke[*config.Config](di),
		authService:   do.MustInvoke[*auth.Service](di),
		pubSubService: do.MustInvoke[*pubsub.Service](di),
	}
}

func (c *WS) handleInternal(conn *websocket.Conn, channels []string) {
	writeChan := make(chan api.IdMessage, 16)
	defer close(writeChan)

	for _, channel := range channels {
		sub := c.pubSubService.Subscribe(channel, func(data any) {
			defer func() {
				if err := recover(); err != nil {
					slog.Warn("Panic in subscription handler", slog.Any("error", err))
				}
			}()

			idMsg, ok := data.(api.IdMessage)
			if !ok {
				slog.LogAttrs(context.Background(), slog.LevelError, "Failed to cast pubsub message to IdMessage",
					slog.Any("data", data),
				)
				return
			}

			writeChan <- idMsg
		})
		defer c.pubSubService.Unsubscribe(sub) // it's ok to defer there
	}

	go func() {
		idCache := ttlcache.New[string, struct{}]()

		go idCache.Start()
		defer idCache.Stop()

		for data := range writeChan {
			id := data.GetId()

			if id != "" && idCache.Has(data.GetId()) {
				continue
			}

			idCache.Set(data.GetId(), struct{}{}, time.Minute)

			_ = conn.SetWriteDeadline(time.Now().Add(1 * time.Minute))
			_ = conn.WriteJSON(data)
		}
	}()

	for {
		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Minute))

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if bytes.Equal(msg, pingMsg) {
			writeChan <- &api.WsMessage{
				Event: "pong",
			}
		}
	}
}

func (c *WS) Handle(conn *websocket.Conn) {
	var channels []string

	if c.authService.IsAdminLocals(conn.Locals) {
		channels = append(channels, "admin")
	}

	c.handleInternal(conn, channels)
}
