package models

import "time"

type ItemDTO struct {
	ItemNO string
	ASIN   string
	Status int
}

type ReviewDTO struct {
	ID           string
	AmazonID     string
	Location     string
	CustomerName string
	Title        string
	Content      string
	IsVerified   bool
	Rating       float32
	CreatedOn    *time.Time
}
