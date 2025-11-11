package dto

import "time"

// BalanceResponse 잔액 응답
type BalanceResponse struct {
	UserID           int64     `json:"user_id"`
	AvailableBalance int64     `json:"available_balance"`
	PendingBalance   int64     `json:"pending_balance"`
	TotalEarned      int64     `json:"total_earned"`
	TotalUsed        int64     `json:"total_used"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TransactionResponse 거래 내역 응답
type TransactionResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Type         string    `json:"type"`
	Amount       int64     `json:"amount"`
	BalanceAfter int64     `json:"balance_after"`
	ReasonType   string    `json:"reason_type"`
	ReasonDetail string    `json:"reason_detail"`
	OrderID      *int64    `json:"order_id,omitempty"`
	EarnedAt     *string   `json:"earned_at,omitempty"`
	ExpiresAt    *string   `json:"expires_at,omitempty"`
	Expired      bool      `json:"expired"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// TransactionsResponse 거래 내역 목록 응답
type TransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}

// ErrorResponse 에러 응답
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

