package campaign

type CampaignService interface {
	GetCampaign(userId int) ([]Campaign, error)
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
