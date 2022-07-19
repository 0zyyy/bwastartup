package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type CampaignService interface {
	GetCampaign(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	AddCampaign(input NewCampaign) (Campaign, error)
	UpdateCampaign(inputId GetCampaignDetailInput, inputData NewCampaign) (Campaign, error)
	SaveCampaignImage(input GetCampaignImageInput, fileLoc string) (CampaignImage, error)
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

func (s *campaignService) UpdateCampaign(inputId GetCampaignDetailInput, inputData NewCampaign) (Campaign, error) {
	campaign, err := s.campaignRepo.FindById(inputId.Id)
	if err != nil {
		return campaign, err
	}
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.campaignRepo.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}
func (s *campaignService) SaveCampaignImage(input GetCampaignImageInput, fileLoc string) (CampaignImage, error) {
	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 0
		_, err := s.campaignRepo.MarkAllAsNon(input.Id)
		if err != nil {
			return CampaignImage{}, err
		}
	}
	campaignImg := CampaignImage{}
	campaignImg.CampaignId = input.Id
	campaignImg.IsPrimary = isPrimary
	campaginImage, err := s.campaignRepo.CreateImage(campaignImg)
	if err != nil {
		return campaginImage, err
	}
	return campaginImage, nil
}
