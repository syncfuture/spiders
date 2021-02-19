package wayfair

import (
	"time"
)

type ItemDTO struct {
	// ItemNo string
	SKU    string
	Status int
}

type ReviewDTO struct {
	// ID           string
	AmazonID     string
	CustomerNo   string
	SKU          string
	Location     string
	CustomerName string
	Title        string
	Content      string
	StripInfo    string
	IsVerified   bool
	Rating       float32
	CreatedOn    *time.Time
}

type ItemQuery struct {
	Cursor   string
	SKU      string
	ItemNo   string
	Status   string
	PageSize int
}

type ItemQueryResult struct {
	MsgCode    string
	Cursor     string
	Items      []*ItemDTO
	TotalCount int64
}

type ReviewQuery struct {
	Cursor   string
	SKU      string
	ItemNo   string
	FromDate string
	PageSize int
}

type ReviewQueryResult struct {
	MsgCode    string
	Cursor     string
	TotalCount int64
	Reviews    []*ReviewDTO
}

type ReviewResult struct {
	Reviews     []*ReviewDTO
	NextPageURL string
}
