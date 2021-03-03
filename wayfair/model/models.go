package model

type ItemDTO struct {
	Items  string
	SKU    string
	URL    string
	Status int
}

// type ReviewDTO struct {
// 	SKU      string
// 	Items    string
// 	Comments string
// 	Rating   string
// 	Date     time.Time
// 	Photos   []string
// 	Name     string
// 	Badge    string
// 	Helpful  int
// }

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

type ReviewResp struct {
	Data *ReviewRespData `json:"data,omitempty"`
}
type ReviewRespData struct {
	Product *ReviewRespProduct `json:"product,omitempty"`
}
type ReviewRespProduct struct {
	CustomerReviews *ReviewRespReviews `json:"customerReviews,omitempty"`
}
type ReviewRespReviews struct {
	Reviews []*ReviewDTO `json:"reviews,omitempty"`
}
type ReviewDTO struct {
	ReviewID                   int            `json:"reviewId,omitempty"`
	ReviewerName               string         `json:"reviewerName,omitempty"`
	HasVerifiedBuyerStatus     bool           `json:"hasVerifiedBuyerStatus,omitempty"`
	IsUSReviewer               bool           `json:"isUSReviewer,omitempty"`
	ReviewerBadgeText          string         `json:"reviewerBadgeText,omitempty"`
	ReviewerBadgeID            int            `json:"reviewerBadgeId,omitempty"`
	RatingStars                int            `json:"ratingStars,omitempty"`
	Date                       string         `json:"date,omitempty"`
	Headline                   string         `json:"headline,omitempty"`
	ProductComments            string         `json:"productComments,omitempty"`
	HeadlineTranslation        string         `json:"headlineTranslation,omitempty"`
	ProductCommentsTranslation string         `json:"productCommentsTranslation,omitempty"`
	LanguageCode               string         `json:"languageCode,omitempty"`
	ReviewHelpful              int            `json:"reviewHelpful,omitempty"`
	IsReviewHelpfulUpvoted     bool           `json:"isReviewHelpfulUpvoted,omitempty"`
	ProductName                string         `json:"productName,omitempty"`
	ProductUrl                 string         `json:"productUrl,omitempty"`
	CustomerPhotos             []*ReviewPhoto `json:"customerPhotos,omitempty"`
	SKU                        string
	Items                      string
}
type ReviewPhoto struct {
	Index     int    `json:"idx,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Src       string `json:"src,omitempty"`
	IreID     int    `json:"ire_id,omitempty"`
}
