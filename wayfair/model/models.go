package model

import (
	"time"
)

type ItemDTO struct {
	Items  string
	SKU    string
	URL    string
	Status int
}

type ReviewDTO struct {
	SKU      string
	Items    string
	Comments string
	Rating   string
	Date     time.Time
	Photos   []string
	Name     string
	Badge    string
	Helpful  int
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
	Reviews []*ReviewDTO
}
