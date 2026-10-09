package tickets

import (
	"time"

	"github.com/yoyojr404/internal/common/entity"
)

type Ticket struct {
	CategoryName           string     `db:"category_name"`
	VerificationID         int64      `db:"verification_id"`
	VerificationVerified   bool       `db:"verification_verified"`
	VerificationVerifiedAt *time.Time `db:"verification_verified_at"`
	entity.Ticket
}
