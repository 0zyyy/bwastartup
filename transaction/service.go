package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"strconv"

	"errors"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(UserID int) ([]Transaction, error)
	AddTransaction(input CreateTransactionInput) (Transaction, error)
	PaymentProcess(payment coreapi.TransactionStatusResponse) error
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	if input.ID <= 0 {
		return []Transaction{}, errors.New("invalid campaign id")
	}

	campaign, err := s.campaignRepository.FindById(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserId != input.User.Id {
		return []Transaction{}, errors.New("not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(UserID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(UserID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) AddTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = input.User.Id
	transaction.Status = "pending"
	transaction.Code = ""
	newTransaction, err := s.repository.Add(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentUrl, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentUrl

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) PaymentProcess(payment coreapi.TransactionStatusResponse) error {
	transactionStatusResp, e := s.paymentService.GetStatus(payment.OrderID)
	if e != nil {
		return e
	} else {
		if transactionStatusResp != nil {
			orderID, err := strconv.Atoi(payment.OrderID)
			if err != nil {
				return err
			}
			transaction, err := s.repository.GetByID(orderID)
			if err != nil {
				return err
			}
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					transaction.Status = "paid"
					_, err = s.repository.Update(transaction)
					if err != nil {
						return err
					}

					campaign, err := s.campaignRepository.FindById(transaction.CampaignID)
					if err != nil {
						return err
					}
					campaign.BackerCount = campaign.BackerCount + 1
					campaign.CurrentAmount = campaign.CurrentAmount + transaction.Amount
					_, err = s.campaignRepository.Update(campaign)
					if err != nil {
						return err
					}

				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				transaction.Status = "paid"
				_, err = s.repository.Update(transaction)
				if err != nil {
					return err
				}
				campaign, err := s.campaignRepository.FindById(transaction.CampaignID)
				if err != nil {
					return err
				}
				campaign.BackerCount = campaign.BackerCount + 1
				campaign.CurrentAmount = campaign.CurrentAmount + transaction.Amount
				_, err = s.campaignRepository.Update(campaign)
				if err != nil {
					return err
				}
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
				transaction.Status = transactionStatusResp.TransactionStatus
				_, err = s.repository.Update(transaction)
				if err != nil {
					return err
				}
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
				transaction.Status = transactionStatusResp.TransactionStatus
				_, err = s.repository.Update(transaction)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
