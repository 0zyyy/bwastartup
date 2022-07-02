package campaign

type CampaignFormatter struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	Name       string `json:"name"`
	ShortDesc  string `json:"short_description"`
	ImageUrl   string `json:"image_url"`
	GoalAmount int    `json:"goal_amount"`
	CurrAmount int    `json:"current_amount"`
}

// singular
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campFormatter := CampaignFormatter{
		Id:         campaign.ID,
		UserId:     campaign.UserId,
		Name:       campaign.Name,
		ShortDesc:  campaign.ShortDesc,
		ImageUrl:   "",
		GoalAmount: campaign.GoalAmount,
		CurrAmount: campaign.CurrAmount,
	}

	if len(campaign.CampaignImages) > 0 {
		campFormatter.ImageUrl = campaign.CampaignImages[0].FileName
	}
	return campFormatter
}

// plural campaign
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		campaignsFormatter = append(campaignsFormatter, FormatCampaign(campaign))
	}
	return campaignsFormatter
}
