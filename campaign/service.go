package campaign

type CampaignService interface {
	GetAllCampaign() ([]Campaign, error)
	GetCampaignByUserId(userId int) ([]Campaign, error)
}

type campaignService struct {
	campaignRepo Repository
}

func NewService(campaignRepo Repository) *campaignService {
	return &campaignService{campaignRepo}
}

func (s *campaignService) GetAllCampaign() ([]Campaign, error) {
	campaings, err := s.campaignRepo.GetCampaigns()
	if err != nil {
		return campaings, err
	}
	return campaings, nil
}
