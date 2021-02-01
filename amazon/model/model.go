package model

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
	SearchAfter string
	Items       []*ItemDTO
	TotalCount  int64
}
