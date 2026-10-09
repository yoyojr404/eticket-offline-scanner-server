package ticketverification

import (
	"context"
	"encoding/json"
	"fmt"
)

type TicketVerificationHandler struct {
	service Service
}

func NewTicketVerificationHandler(service Service) *TicketVerificationHandler {
	return &TicketVerificationHandler{service: service}
}

func (h *TicketVerificationHandler) Handle(ctx context.Context, msg json.RawMessage) ([]byte, error) {
	var payload Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		return nil, err
	}

	if err := h.service.VerifyTicket(ctx, payload.ID); err != nil {
		return nil, err
	}

	return json.Marshal(fmt.Sprintf("Ticket #%d verified", payload.ID))
}
