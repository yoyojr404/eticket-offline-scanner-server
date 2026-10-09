package tickets

type TicketResult struct {
	Entities     []Ticket
	CurrentPage  int
	TotalPerPage int
	Total        int64
}
