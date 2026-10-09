package ws

import (
	"log"
	"sync"
	"time"

	"github.com/yoyojr404/internal/ws/handler"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan Response
	receive   chan handler.Message
	closeOnce sync.Once
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan Response, 256),
		receive:   make(chan handler.Message, 10),
		closeOnce: sync.Once{},
	}
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		c.hub.unregister <- c
	})
}

func (c *Client) ReadPump() {
	defer func() {
		c.close()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg handler.Message
		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}

			break
		}

		select {
		case c.receive <- msg:
		default:
			log.Printf("receive channel full, dropping message")
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.close()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)

			if err != nil {
				log.Printf("write error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping error: %v", err)
				return
			}
		}
	}
}
