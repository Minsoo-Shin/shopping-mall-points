package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shopping-mall/internal/handler/dto"
	pointUseCase "shopping-mall/internal/usecase/point"

	"github.com/gorilla/mux"
)

// OrderHandler 주문 핸들러 (포인트 관련)
type OrderHandler struct {
	useUseCase    *pointUseCase.UsePointsUseCase
	earnUseCase   *pointUseCase.EarnPointsUseCase
	refundUseCase *pointUseCase.RefundPointsUseCase
}

// NewOrderHandler 주문 핸들러 생성
func NewOrderHandler(
	useUseCase *pointUseCase.UsePointsUseCase,
	earnUseCase *pointUseCase.EarnPointsUseCase,
	refundUseCase *pointUseCase.RefundPointsUseCase,
) *OrderHandler {
	return &OrderHandler{
		useUseCase:    useUseCase,
		earnUseCase:   earnUseCase,
		refundUseCase: refundUseCase,
	}
}

// ConfirmOrder 주문 확정 (포인트 적립)
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid order_id")
		return
	}

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

	req.OrderID = orderID
	ctx := r.Context()
	if err := h.earnUseCase.EarnPointsFromPurchase(ctx, userID, req.PaymentAmount, req.OrderID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "order confirmed and points earned"})
}

// RefundOrder 주문 환불 (포인트 복구/회수)
func (h *OrderHandler) RefundOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid order_id")
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()
	if err := h.refundUseCase.RefundPoints(ctx, userID, orderID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "order refunded and points processed"})
}
