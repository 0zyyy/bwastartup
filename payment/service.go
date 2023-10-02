package payment

import (
	"bwastartup/configs"
	"bwastartup/user"
	"log"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
	client coreapi.Client
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	GetStatus(orderID string) (*coreapi.TransactionStatusResponse, error)
}

func setupGlobalMidtransConfig() {
	midtrans.ServerKey = configs.SandboxServerKey
	midtrans.Environment = midtrans.Sandbox
}

func NewService(client coreapi.Client) *service {
	setupGlobalMidtransConfig()
	return &service{client: client}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	responseToken, err := snap.CreateTransaction(snapReq)
	if err != nil {
		log.Println("error", err)
		return "", err
	}
	log.Println("responseToken", responseToken)
	return responseToken.RedirectURL, nil
}

func (s *service) GetStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	response, err := s.client.CheckTransaction(orderID)
	if err != nil {
		return &coreapi.TransactionStatusResponse{}, err
	}
	log.Println("response", response)
	return response, nil
}
