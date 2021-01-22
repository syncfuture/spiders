package items

import "time"

type ItemDTO struct {
	ItemNO string
	ASIN   string
	Status int
}

type ReviewDTO struct {
	ID           string
	Location     string
	CustomerName string
	Title        string
	Content      string
	IsVerified   bool
	Rating       float32
	CreatedOn    *time.Time
}
