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
	ASIN     string
	ItemNo   string
	From     string
	To       string
	PageSize int
}

type ReviewQueryResult struct {
	MsgCode    string
	Cursor     string
	TotalCount int64
	Reviews    []*amazon.ReviewDTO
}