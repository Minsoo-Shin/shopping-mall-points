package point

import "time"

// TransactionType 거래 유형
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "EARN"   // 적립
	TransactionTypeUse    TransactionType = "USE"    // 사용
	TransactionTypeExpire TransactionType = "EXPIRE" // 만료
	TransactionTypeCancel TransactionType = "CANCEL" // 취소
)

// ReasonType 적립/사용 사유
type ReasonType string

const (
	ReasonTypePurchase ReasonType = "PURCHASE" // 구매
	ReasonTypeReview   ReasonType = "REVIEW"   // 리뷰
	ReasonTypeSignup   ReasonType = "SIGNUP"   // 가입
	ReasonTypeRefund   ReasonType = "REFUND"   // 환불
	ReasonTypeAdmin    ReasonType = "ADMIN"    // 관리자
)

// TransactionStatus 거래 상태
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"   // 대기
	TransactionStatusConfirmed TransactionStatus = "CONFIRMED" // 확정
	TransactionStatusCancelled TransactionStatus = "CANCELLED" // 취소
)

// Transaction 포인트 거래 내역
type Transaction struct {
	ID            int64
	UserID        int64
	Type          TransactionType
	Amount        int64
	BalanceAfter  int64
	ReasonType    ReasonType
	ReasonDetail  string
	OrderID       *int64
	EarnedAt      *time.Time
	ExpiresAt     *time.Time
	Expired       bool
	Status        TransactionStatus
	CreatedAt     time.Time
}

// IsExpired 만료 여부 확인
func (t *Transaction) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return t.Expired || time.Now().After(*t.ExpiresAt)
}

