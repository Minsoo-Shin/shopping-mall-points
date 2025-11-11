package point

import "time"

// Policy 포인트 정책
type Policy struct {
	EarnRate          float64 // 적립률 (0.05 = 5%)
	ReviewTextPoints  int64   // 텍스트 리뷰 적립
	ReviewPhotoPoints int64   // 포토 리뷰 적립
	SignupBonus       int64   // 가입 보너스
	MinOrderAmount    int64   // 최소 주문 금액
	MaxEarnPerOrder   int64   // 주문당 최대 적립
	ExpiryMonths      int     // 유효기간 (월)
	EarnDelayDays     int     // 적립 지연 일수
	
	MinUseAmount      int64   // 최소 사용 금액
	UseUnit           int64   // 사용 단위
	MaxUseRate        float64 // 최대 사용 비율 (0.5 = 50%)
	MinPaymentAmount  int64   // 최소 결제 금액
}

// NewDefaultPolicy 기본 정책 생성
func NewDefaultPolicy() *Policy {
	return &Policy{
		EarnRate:          0.05,  // 5%
		ReviewTextPoints:  100,
		ReviewPhotoPoints: 500,
		SignupBonus:       3000,
		MinOrderAmount:    10000,
		MaxEarnPerOrder:   50000,
		ExpiryMonths:      12,
		EarnDelayDays:     7,
		
		MinUseAmount:     1000,
		UseUnit:          100,
		MaxUseRate:       0.5, // 50%
		MinPaymentAmount: 1000,
	}
}

// CalculateEarnPoints 적립 포인트 계산
func (p *Policy) CalculateEarnPoints(paymentAmount int64) int64 {
	earnPoints := int64(float64(paymentAmount) * p.EarnRate)
	if earnPoints > p.MaxEarnPerOrder {
		return p.MaxEarnPerOrder
	}
	return earnPoints
}

// ValidateUse 사용 유효성 검증
func (p *Policy) ValidateUse(useAmount, orderAmount, availableBalance int64) error {
	// 최소 사용 금액 체크
	if useAmount < p.MinUseAmount {
		return ErrBelowMinUseAmount
	}
	
	// 사용 단위 체크
	if useAmount%p.UseUnit != 0 {
		return ErrInvalidUseUnit
	}
	
	// 보유 포인트 확인
	if availableBalance < useAmount {
		return ErrInsufficientPoints
	}
	
	// 최대 사용 비율 체크
	maxUseAmount := int64(float64(orderAmount) * p.MaxUseRate)
	if useAmount > maxUseAmount {
		return ErrExceedMaxUseRate
	}
	
	// 최소 결제 금액 체크 (전액 포인트 결제 방지)
	paymentAmount := orderAmount - useAmount
	if paymentAmount < p.MinPaymentAmount {
		return ErrBelowMinPayment
	}
	
	return nil
}

// CalculateExpiryDate 만료일 계산
func (p *Policy) CalculateExpiryDate(earnedAt time.Time) time.Time {
	return earnedAt.AddDate(0, p.ExpiryMonths, 0)
}

// CalculateEarnDate 실제 적립일 계산 (구매 확정 후 지연 일수)
func (p *Policy) CalculateEarnDate(confirmedAt time.Time) time.Time {
	return confirmedAt.AddDate(0, 0, p.EarnDelayDays)
}

