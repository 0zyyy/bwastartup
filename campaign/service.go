package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type CampaignService interface {
	GetCampaign(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	AddCampaign(input NewCampaign) (Campaign, error)
}

type campaignService struct {
	campaignRepo Repository
}

func NewService(campaignRepo Repository) *campaignService {
	return &campaignService{campaignRepo}
}

func (s *campaignService) GetCampaign(userId int) ([]Campaign, error) {
	if userId != 0 {
		campaign, err := s.campaignRepo.FindByUserId(userId)
		if err != nil {
			return campaign, err
		}
		return campaign, nil
	}
	campaign, err := s.campaignRepo.FindAll()
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignService) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.campaignRepo.FindById(input.Id)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignService) AddCampaign(input NewCampaign) (Campaign, error) {
	newCampaign := Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		GoalAmount:       input.GoalAmount,
		Description:      input.Description,
		Perks:            input.Perks,
		UserId:           input.User.Id,
	}

	// pembuatan slug
	slugCand := fmt.Sprintf("%s %d", input.Name, input.User.Id)
	newCampaign.Slug = slug.Make(slugCand)
	newCampaign, err := s.campaignRepo.CreateCampaign(newCampaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}
