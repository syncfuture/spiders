package model

import "github.com/syncfuture/scraper/amazon"

type ItemDTO struct {
	ASIN   string
	ItemNo string
	Status int
}

type ItemQuery struct {
	Cursor   string
	ASIN     string
	ItemNo   string
	PageSize int
	Status   int
}

type ItemQueryResult struct {
	MsgCode    string
	Cursor     string
	Items      []*ItemDTO
	TotalCount int64
}

type ReviewQueryResult struct {
	MsgCode    string
	Cursor     string
	TotalCount int64
	Reviews    []*amazon.ReviewDTO
}
