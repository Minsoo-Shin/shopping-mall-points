package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	pointDomain "shopping-mall/internal/domain/point"
	"shopping-mall/internal/handler/dto"
	pointUseCase "shopping-mall/internal/usecase/point"

	"github.com/gorilla/mux"
)

// PointHandler 포인트 핸들러
type PointHandler struct {
	queryUseCase *pointUseCase.QueryPointsUseCase
	useUseCase   *pointUseCase.UsePointsUseCase
	earnUseCase  *pointUseCase.EarnPointsUseCase
}

// NewPointHandler 포인트 핸들러 생성
func NewPointHandler(
	queryUseCase *pointUseCase.QueryPointsUseCase,
	useUseCase *pointUseCase.UsePointsUseCase,
	earnUseCase *pointUseCase.EarnPointsUseCase,
) *PointHandler {
	return &PointHandler{
		queryUseCase: queryUseCase,
		useUseCase:   useUseCase,
		earnUseCase:  earnUseCase,
	}
}

// GetBalance 잔액 조회
func (h *PointHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()
	userPoint, err := h.queryUseCase.GetBalance(ctx, userID)
	if err != nil {
		if err == pointDomain.ErrPointNotFound {
			respondError(w, http.StatusNotFound, "point not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.BalanceResponse{
		UserID:           userPoint.UserID,
		AvailableBalance: userPoint.AvailableBalance,
		PendingBalance:   userPoint.PendingBalance,
		TotalEarned:      userPoint.TotalEarned,
		TotalUsed:        userPoint.TotalUsed,
		UpdatedAt:        userPoint.UpdatedAt,
	})
}

// GetTransactions 거래 내역 조회
func (h *PointHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	limit, offset := getPagination(r)

	ctx := r.Context()
	transactions, err := h.queryUseCase.GetTransactions(ctx, userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = toTransactionResponse(tx)
	}

	respondJSON(w, http.StatusOK, dto.TransactionsResponse{
		Transactions: responses,
		Total:        len(responses),
		Limit:        limit,
		Offset:       offset,
	})
}

// UsePoints 포인트 사용
func (h *PointHandler) UsePoints(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	var req dto.UsePointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err := h.useUseCase.UsePoints(ctx, userID, req.UseAmount, req.OrderAmount, req.OrderID); err != nil {
		switch err {
		case pointDomain.ErrInsufficientPoints:
			respondError(w, http.StatusBadRequest, "insufficient points")
		case pointDomain.ErrBelowMinUseAmount:
			respondError(w, http.StatusBadRequest, "below minimum use amount")
		case pointDomain.ErrInvalidUseUnit:
			respondError(w, http.StatusBadRequest, "invalid use unit")
		case pointDomain.ErrExceedMaxUseRate:
			respondError(w, http.StatusBadRequest, "exceed maximum use rate")
		case pointDomain.ErrBelowMinPayment:
			respondError(w, http.StatusBadRequest, "below minimum payment amount")
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "points used successfully"})
}

// EarnPoints 포인트 적립
func (h *PointHandler) EarnPoints(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	var req dto.EarnPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err := h.earnUseCase.EarnPointsFromPurchase(ctx, userID, req.PaymentAmount, req.OrderID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "points earned successfully"})
}

// Helper functions

func getUserID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		// Query parameter에서 시도
		userIDStr = r.URL.Query().Get("user_id")
		if userIDStr == "" {
			return 0, http.ErrMissingFile
		}
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func getPagination(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20 // 기본값
	offset = 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}

func toTransactionResponse(tx *pointDomain.Transaction) dto.TransactionResponse {
	resp := dto.TransactionResponse{
		ID:           tx.ID,
		UserID:       tx.UserID,
		Type:         string(tx.Type),
		Amount:       tx.Amount,
		BalanceAfter: tx.BalanceAfter,
		ReasonType:   string(tx.ReasonType),
		ReasonDetail: tx.ReasonDetail,
		OrderID:      tx.OrderID,
		Expired:      tx.Expired,
		Status:       string(tx.Status),
		CreatedAt:    tx.CreatedAt,
	}

	if tx.EarnedAt != nil {
		earnedAtStr := tx.EarnedAt.Format(time.RFC3339)
		resp.EarnedAt = &earnedAtStr
	}

	if tx.ExpiresAt != nil {
		expiresAtStr := tx.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &expiresAtStr
	}

	return resp
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, dto.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
