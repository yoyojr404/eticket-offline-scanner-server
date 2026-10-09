package ws

import (
	"context"
	"sync"

	"github.com/yoyojr404/internal/ws/handler"
)

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	dispatcher *handler.Dispatcher
}

func NewHub(dispatcher *handler.Dispatcher) *Hub {
	return &Hub{
		mu:         sync.RWMutex{},
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		dispatcher: dispatcher,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client != nil {
				h.clients[client] = true
			}
		case client := <-h.unregister:
			if client != nil {
				h.mu.Lock()
				if h.clients[client] {
					delete(h.clients, client)
					close(client.receive)
					close(client.send)
				}
				h.mu.Unlock()
			}
		}
	}
}

func (h *Hub) SendTo(client *Client, message Response) bool {
	if client == nil {
		return false
	}

	h.mu.RLock()
	_, ok := h.clients[client]
	h.mu.RUnlock()
	if !ok {
		return false
	}

	select {
	case client.send <- message:
		return true
	default:
		return false
	}
}

func (h *Hub) HandleMessage(ctx context.Context, client *Client) {
	defer func() {
		client.close()
		client.conn.Close()
	}()

	for msg := range client.receive {
		result, err := h.dispatcher.Dispatch(ctx, msg)

		response := Response{}
		if err != nil {
			Error(err, &response)
		} else {
			Success(result, &response)
		}

		h.SendTo(client, response)
	}
}
