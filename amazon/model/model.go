package model

import "github.com/syncfuture/scraper/amazon"

type ItemDTO struct {
	ASIN        string
	ItemNo      string
	SearchAfter string
	Status      int
}

type ItemQuery struct {
	SearchAfter string
	ASIN        string
	ItemNo      string
	PageSize    int
	Status      int
}

type ItemQueryResult struct {
	MsgCode     string
	SearchAfter string
	Items       []*ItemDTO
	TotalCount  int64
}

type ReviewQueryResult struct {
	MsgCode     string
	SearchAfter string
	Reviews     []*amazon.ReviewDTO
	TotalCount  int64
}
