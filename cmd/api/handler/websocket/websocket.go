package websocket

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Manager struct {
	upgrader websocket.Upgrader
	clientId *atomic.Uint64
	clients  Map[uint64, *client]
}

func NewManager() *Manager {
	return &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clientId: new(atomic.Uint64),
		clients:  NewMap[uint64, *client](),
	}
}

func (manager *Manager) Handle(c echo.Context) error {
	ws, err := manager.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	ws.SetReadLimit(1024 * 10) // 10KB

	sId := manager.clientId.Add(1)
	sub := newClient(sId, ws, manager)

	manager.clients.Set(sId, sub)

	go sub.ReadMessages(c.Request().Context(), ws, sub, c.Logger())
	go sub.WriteMessages(c.Request().Context(), c.Logger())

	c.Logger().Infof("client %d connected", sId)
	return nil
}

func (manager *Manager) Close() error {
	return manager.clients.Range(func(_ uint64, value *client) (error, bool) {
		if err := value.Close(); err != nil {
			return err, false
		}
		return nil, false
	})
}

func (manager *Manager) NotifyAll(msg any) {
	_ = manager.clients.Range(func(_ uint64, value *client) (error, bool) {
		value.Notify(msg)
		return nil, false
	})
}
