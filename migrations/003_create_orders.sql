-- orders 테이블 생성 (포인트 시스템과 연관)
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '사용자 ID',
    total_amount BIGINT NOT NULL COMMENT '총 주문 금액',
    point_used BIGINT NOT NULL DEFAULT 0 COMMENT '사용한 포인트',
    point_to_earn BIGINT NOT NULL DEFAULT 0 COMMENT '적립 예정 포인트',
    payment_amount BIGINT NOT NULL COMMENT '실제 결제 금액',
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' COMMENT '주문 상태',
    confirmed_at TIMESTAMP NULL COMMENT '구매 확정 시점',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='주문';

