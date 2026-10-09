package entity

type Ticket struct {
	ID         int64  `db:"id"`
	Code       string `db:"code"`
	CategoryID int64  `db:"category_id"`
}
