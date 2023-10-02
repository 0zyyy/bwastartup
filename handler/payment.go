package handler

import (
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/coreapi"
)

type paymentHandler struct {
	service            payment.Service
	transactionService transaction.Service
}

func NewPaymentHandler(service payment.Service, transactionService transaction.Service) *paymentHandler {
	return &paymentHandler{service, transactionService}
}

func (h *paymentHandler) Notification(c *gin.Context) {
	// 1. Initialize empty map
	var input coreapi.TransactionStatusResponse
	// 2. Parse JSON request body and use it to set json to payload
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Error to get notification", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 4. Check transaction to Midtrans with param orderId
	err = h.transactionService.PaymentProcess(input)
	if err != nil {
		response := helper.APIResponse("Error to process notification", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
