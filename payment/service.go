package payment

import (
	"bwastartup/configs"
	"bwastartup/user"
	"crypto/sha512"
	"encoding/hex"
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
	CheckSignature(signatureInput SignatureInput) (bool, error)
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
	return responseToken.RedirectURL, nil
}

func (s *service) GetStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	response, err := s.client.CheckTransaction(orderID)
	if err != nil {
		return &coreapi.TransactionStatusResponse{}, err
	}
	log.Println("response", response.SignatureKey)
	return response, nil
}

func (s *service) CheckSignature(signatureInput SignatureInput) (bool, error) {
	signature := sha512.New()
	signature.Write([]byte(signatureInput.OrderID + signatureInput.StatusCode + signatureInput.GrossAmount + configs.SandboxServerKey))
	signatureKey := hex.EncodeToString(signature.Sum(nil))
	if signatureKey != signatureInput.SignatureKey {
		return false, nil
	}
	return true, nil
}
