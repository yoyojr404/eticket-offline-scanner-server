package handler

import (
	"context"
	"fmt"
)

type Dispatcher struct {
	handlers map[MessageType]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[MessageType]Handler),
	}
}

func (d *Dispatcher) Register(msgType MessageType, h Handler) {
	d.handlers[msgType] = h
}

func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) ([]byte, error) {
	h, ok := d.handlers[msg.Type]
	if !ok {
		return nil, fmt.Errorf("handler not found for message type: %v", msg.Type)
	}
	return h.Handle(ctx, msg.Payload)
}
