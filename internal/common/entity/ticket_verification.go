package entity

import "time"

type TicketVerification struct {
	ID         int64      `db:"id"`
	TicketID   int64      `db:"ticket_id"`
	Verified   bool       `db:"verified"`
	CreatedAt  time.Time  `db:"created_at"`
	VerifiedAt *time.Time `db:"verified_at"`
}
