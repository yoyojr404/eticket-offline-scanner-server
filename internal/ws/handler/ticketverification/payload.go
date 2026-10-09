package ticketverification

type Payload struct {
	ID       int64 `json:"id"`
	Verified bool  `json:"verified"`
}
