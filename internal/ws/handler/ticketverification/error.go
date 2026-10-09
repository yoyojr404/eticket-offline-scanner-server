package ticketverification

import "errors"

var (
	ErrTicketVerificationNotFound        = errors.New("ticket verification not found")
	ErrTicketVerificationHasBeenVerified = errors.New("ticket verification has been verified")
)
