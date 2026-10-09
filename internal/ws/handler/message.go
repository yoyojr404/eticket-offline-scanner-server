package handler

import "encoding/json"

type MessageType int

const (
	CmdHello        MessageType = 0
	CmdVerifyTicket MessageType = 1
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
