package dto

// UsePointsRequest 포인트 사용 요청
type UsePointsRequest struct {
	OrderID     int64 `json:"order_id"`
	UseAmount   int64 `json:"use_amount"`
	OrderAmount int64 `json:"order_amount"`
}

// EarnPointsRequest 포인트 적립 요청
type EarnPointsRequest struct {
	OrderID       int64 `json:"order_id"`
	PaymentAmount int64 `json:"payment_amount"`
}

// ReviewPointsRequest 리뷰 포인트 적립 요청
type ReviewPointsRequest struct {
	IsPhoto bool `json:"is_photo"`
}

