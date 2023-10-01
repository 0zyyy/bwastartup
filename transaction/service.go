package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(UserID int) ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
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
