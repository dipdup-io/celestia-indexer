package websocket

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type client struct {
	id      uint64
	ws      *websocket.Conn
	manager *Manager
	filters *filters
	ch      chan any
	close   chan struct{}
	wg      *sync.WaitGroup
}

type filters struct {
	head bool
}

func newClient(id uint64, ws *websocket.Conn, manager *Manager) *client {
	return &client{
		id:      id,
		ws:      ws,
		manager: manager,
		filters: new(filters),
		ch:      make(chan any, 1024),
		close:   make(chan struct{}),
		wg:      new(sync.WaitGroup),
	}
}

func (c *client) Add(msg Subscribe) error {
	switch msg.Channel {
	case ChannelHead:
		c.filters.head = true
	default:
		return errors.Wrap(ErrUnknownChannel, msg.Channel)
	}
	return nil
}

func (c *client) Remove(msg Unsubscribe) error {
	switch msg.Channel {
	case ChannelHead:
		c.filters.head = false
	default:
		return errors.Wrap(ErrUnknownChannel, msg.Channel)
	}
	return nil
}

func (c *client) Notify(msg any) {
	c.ch <- msg
}

func (c *client) Close() error {
	c.wg.Wait()

	if err := c.ws.Close(); err != nil {
		return err
	}

	close(c.close)
	close(c.ch)
	return nil
}

func (c *client) WriteMessages(ctx context.Context, log echo.Logger) {
	c.wg.Add(1)
	defer c.wg.Done()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-c.close:
			return

		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("writemsg: %s", err)
				return
			}

		case msg, ok := <-c.ch:
			if !ok {
				if err := c.ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Errorf("send close message: %s", err)
				}
				return
			}

			switch typ := msg.(type) {
			case storage.Block:
				if c.filters.head {
					if err := c.ws.WriteJSON(typ); err != nil {
						log.Errorf("send head: %s", err)
					}
				}
			default:
				log.Errorf("unknown message type from notification: %T", msg)
			}
		}
	}
}

func (c *client) ReadMessages(ctx context.Context, ws *websocket.Conn, sub *client, log echo.Logger) {
	c.wg.Add(1)
	defer func() {
		c.manager.clients.Delete(sub.id)
		c.wg.Done()
		c.close <- struct{}{}
		c.Close()
	}()

	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error(err)
		return
	}
	c.ws.SetPongHandler(c.pongHandler)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.read(ctx, ws); err != nil {
				switch {
				case err == io.EOF:
					return
				case err == ErrTimeout:
					return
				case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
					return
				}
				log.Errorf("read websocket message: %s", err.Error())
			}
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *client) pongHandler(pongMsg string) error {
	return c.ws.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *client) read(ctx context.Context, ws *websocket.Conn) error {
	var msg Message
	if err := c.ws.ReadJSON(&msg); err != nil {
		return err
	}

	switch msg.Method {
	case MethodSubscribe:
		return c.handleSubscribeMessage(ctx, msg)
	case MethodUnsubscribe:
		return c.handleUnsubscribeMessage(ctx, msg)
	default:
		return errors.Wrap(ErrUnknownMethod, msg.Method)
	}
}

func (c *client) handleSubscribeMessage(ctx context.Context, msg Message) error {
	var subscribeMsg Subscribe
	if err := json.UnmarshalContext(ctx, msg.Body, &subscribeMsg); err != nil {
		return err
	}

	return c.Add(subscribeMsg)
}

func (c *client) handleUnsubscribeMessage(ctx context.Context, msg Message) error {
	var unsubscribeMsg Unsubscribe
	if err := json.UnmarshalContext(ctx, msg.Body, &unsubscribeMsg); err != nil {
		return err
	}
	return c.Remove(unsubscribeMsg)
}
