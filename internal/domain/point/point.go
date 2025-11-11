package point

import "time"

// UserPoint 사용자 포인트 집계 루트
type UserPoint struct {
	UserID          int64
	AvailableBalance int64 // 사용 가능 포인트
	PendingBalance   int64 // 적립 예정 포인트
	TotalEarned      int64 // 누적 적립
	TotalUsed        int64 // 누적 사용
	UpdatedAt        time.Time
}

// CanUse 사용 가능 여부 확인
func (up *UserPoint) CanUse(amount int64) error {
	if up.AvailableBalance < amount {
		return ErrInsufficientPoints
	}
	return nil
}

// Use 포인트 사용
func (up *UserPoint) Use(amount int64) error {
	if err := up.CanUse(amount); err != nil {
		return err
	}
	up.AvailableBalance -= amount
	up.TotalUsed += amount
	up.UpdatedAt = time.Now()
	return nil
}

// Earn 포인트 적립
func (up *UserPoint) Earn(amount int64) {
	up.AvailableBalance += amount
	up.TotalEarned += amount
	up.UpdatedAt = time.Now()
}

// AddPending 적립 예정 포인트 추가
func (up *UserPoint) AddPending(amount int64) {
	up.PendingBalance += amount
	up.UpdatedAt = time.Now()
}

// ConfirmPending 적립 예정 포인트를 실제 적립으로 전환
func (up *UserPoint) ConfirmPending(amount int64) {
	if up.PendingBalance >= amount {
		up.PendingBalance -= amount
		up.AvailableBalance += amount
		up.TotalEarned += amount
		up.UpdatedAt = time.Now()
	}
}

// Refund 포인트 환불
func (up *UserPoint) Refund(amount int64) {
	up.AvailableBalance += amount
	up.UpdatedAt = time.Now()
}

// Expire 포인트 만료
func (up *UserPoint) Expire(amount int64) {
	if up.AvailableBalance >= amount {
		up.AvailableBalance -= amount
		up.UpdatedAt = time.Now()
	}
}

