# 쇼핑몰 포인트 시스템

Go 언어 기반 Clean Architecture로 구현된 포인트 적립/사용/만료 시스템입니다.

## 주요 기능

- 포인트 적립 (구매, 리뷰, 가입 보너스)
- 포인트 사용 (FIFO 방식)
- 포인트 환불
- 포인트 만료 (배치 처리)
- 거래 내역 조회

## 프로젝트 구조

```
shopping-mall/
├── cmd/
│   ├── api/main.go              # API 서버
│   └── worker/main.go           # 배치 작업 (포인트 만료)
├── internal/
│   ├── domain/                  # 도메인 모델 & 비즈니스 로직
│   ├── usecase/                 # 유스케이스
│   ├── repository/              # 데이터 액세스 구현
│   ├── handler/                 # HTTP 핸들러
│   └── infrastructure/          # 외부 의존성
├── pkg/                         # Public 패키지
├── config/                      # 설정
└── migrations/                  # DB 마이그레이션
```

## 설정

환경 변수를 통해 설정할 수 있습니다:

```bash
# 서버 설정
export SERVER_PORT=8080
export ENV=development

# MySQL 설정 (필수)
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password  # MySQL root 비밀번호
export MYSQL_DATABASE=shopping_mall

# Redis 설정 (선택사항)
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0
```

### 빠른 시작 (Docker 사용 시)

```bash
# 1. MySQL Docker 컨테이너 실행
docker run --name mysql-shopping-mall \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=shopping_mall \
  -p 3306:3306 \
  -d mysql:8.0

# 2. 환경 변수 설정
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=password
export MYSQL_DATABASE=shopping_mall

# 3. 마이그레이션 실행
go run cmd/migrate/main.go

# 4. API 서버 실행
go run cmd/api/main.go
```

## 데이터베이스 초기화

애플리케이션은 시작 시 자동으로 데이터베이스를 초기화합니다:

1. **자동 데이터베이스 생성**: 지정한 데이터베이스가 없으면 자동으로 생성합니다.
2. **자동 마이그레이션**: `migrations/` 디렉토리의 SQL 파일들을 자동으로 실행합니다.

### 수동 마이그레이션 실행

마이그레이션을 수동으로 실행하려면:

```bash
# 마이그레이션 도구 사용
go run cmd/migrate/main.go

# 또는 빌드 후 실행
go build -o bin/migrate ./cmd/migrate
./bin/migrate
```

### MySQL CLI로 직접 실행

```bash
mysql -u root -p shopping_mall < migrations/001_create_user_points.sql
mysql -u root -p shopping_mall < migrations/002_create_point_transactions.sql
mysql -u root -p shopping_mall < migrations/003_create_orders.sql
```

## 실행

### API 서버 실행

```bash
go run cmd/api/main.go
```

### Worker 실행 (포인트 만료 배치)

```bash
go run cmd/worker/main.go
```

## API 엔드포인트

### 포인트 조회
- `GET /api/v1/points/balance?user_id={user_id}` - 잔액 조회
- `GET /api/v1/points/transactions?user_id={user_id}&limit={limit}&offset={offset}` - 거래 내역 조회

### 포인트 사용/적립
- `POST /api/v1/points/use` - 포인트 사용
- `POST /api/v1/points/earn` - 포인트 적립

### 주문 관련
- `POST /api/v1/orders/{id}/confirm` - 주문 확정 (포인트 적립)
- `POST /api/v1/orders/{id}/refund` - 주문 환불 (포인트 복구/회수)

## 포인트 정책

### 적립 정책
- 구매 적립률: 결제 금액의 5%
- 리뷰 적립: 텍스트 100P, 포토 500P
- 가입 보너스: 3,000P
- 최소 주문 금액: 10,000원 이상
- 주문당 최대 적립: 50,000P
- 유효기간: 적립일로부터 12개월

### 사용 정책
- 최소 사용: 1,000원 이상
- 사용 단위: 100원 단위
- 최대 사용 비율: 주문 금액의 50%
- 최소 결제 금액: 1,000원 이상 (전액 포인트 결제 방지)
- 차감 방식: FIFO (만료일이 가까운 순서대로)

## 기술 스택

- Go 1.21+
- MySQL
- Redis
- Gorilla Mux
- Zap (로깅)

## 아키텍처

Clean Architecture 원칙을 따르며, 다음과 같은 계층 구조를 가집니다:

- **Domain**: 비즈니스 로직과 도메인 모델
- **UseCase**: 애플리케이션 로직
- **Repository**: 데이터 액세스 추상화
- **Handler**: HTTP 요청/응답 처리
- **Infrastructure**: 외부 의존성 (DB, Cache, Logger)

