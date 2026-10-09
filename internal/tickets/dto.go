package tickets

import "time"

// TicketResponse defines the output representation of a ticket.
type TicketResponse struct {
	ID             int64      `json:"id"`
	Code           string     `json:"code"`
	CategoryID     int64      `json:"category_id"`
	CategoryName   string     `json:"category_name"`
	VerificationID int64      `json:"verification_id"`
	Verified       bool       `json:"verified"`
	VerifiedAt     *time.Time `json:"verified_at"`
}

// ToResponse maps a Ticket entity to TicketResponse DTO.
func ToResponse(entity *Ticket) TicketResponse {
	return TicketResponse{
		ID:             entity.ID,
		Code:           entity.Code,
		CategoryID:     entity.CategoryID,
		CategoryName:   entity.CategoryName,
		VerificationID: entity.VerificationID,
		Verified:       entity.VerificationVerified,
		VerifiedAt:     entity.VerificationVerifiedAt,
	}
}

// ToResponseList maps a slice of Ticket entities to a slice of TicketResponse DTOs.
func ToResponseList(entities []Ticket) []TicketResponse {
	res := make([]TicketResponse, len(entities))
	for i, entity := range entities {
		res[i] = ToResponse(&entity)
	}
	return res
}
