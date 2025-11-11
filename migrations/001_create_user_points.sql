-- user_points 테이블 생성
CREATE TABLE IF NOT EXISTS user_points (
    user_id BIGINT PRIMARY KEY,
    available_balance BIGINT NOT NULL DEFAULT 0 COMMENT '사용 가능 포인트',
    pending_balance BIGINT NOT NULL DEFAULT 0 COMMENT '적립 예정 포인트',
    total_earned BIGINT NOT NULL DEFAULT 0 COMMENT '누적 적립',
    total_used BIGINT NOT NULL DEFAULT 0 COMMENT '누적 사용',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='사용자 포인트 잔액';

