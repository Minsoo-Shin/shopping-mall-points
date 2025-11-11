-- point_transactions 테이블 생성
CREATE TABLE IF NOT EXISTS point_transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '사용자 ID',
    transaction_type ENUM('EARN', 'USE', 'EXPIRE', 'CANCEL') NOT NULL COMMENT '거래 유형',
    amount BIGINT NOT NULL COMMENT '거래 금액',
    balance_after BIGINT NOT NULL COMMENT '거래 후 잔액',
    reason_type ENUM('PURCHASE', 'REVIEW', 'SIGNUP', 'REFUND', 'ADMIN') NOT NULL COMMENT '적립/사용 사유',
    reason_detail VARCHAR(255) DEFAULT '' COMMENT '상세 사유',
    order_id BIGINT NULL COMMENT '주문 ID',
    earned_at TIMESTAMP NULL COMMENT '적립 시점',
    expires_at TIMESTAMP NULL COMMENT '만료 예정일',
    expired BOOLEAN NOT NULL DEFAULT FALSE COMMENT '만료 여부',
    status ENUM('PENDING', 'CONFIRMED', 'CANCELLED') NOT NULL DEFAULT 'PENDING' COMMENT '거래 상태',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id),
    INDEX idx_transaction_type (transaction_type),
    INDEX idx_expires_at (expires_at),
    INDEX idx_created_at (created_at),
    INDEX idx_user_type_status (user_id, transaction_type, status),
    FOREIGN KEY (user_id) REFERENCES user_points(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='포인트 거래 내역';

