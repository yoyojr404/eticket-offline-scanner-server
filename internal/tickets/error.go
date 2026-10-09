package tickets

import "errors"

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrDuplicateCode  = errors.New("ticket code already exists")
)
