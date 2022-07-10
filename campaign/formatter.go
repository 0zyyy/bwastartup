package campaign

import "strings"

type CampaignFormatter struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	Name       string `json:"name"`
	ShortDesc  string `json:"short_description"`
	ImageUrl   string `json:"image_url"`
	GoalAmount int    `json:"goal_amount"`
	CurrAmount int    `json:"current_amount"`
}
type DetailCampaignFormatter struct {
	Id         int            `json:"id"`
	UserId     int            `json:"user_id"`
	Name       string         `json:"name"`
	ShortDesc  string         `json:"short_description"`
	ImageUrl   string         `json:"image_url"`
	GoalAmount int            `json:"goal_amount"`
	CurrAmount int            `json:"current_amount"`
	Perks      []string       `json:"perks"`
	User       UserDetail     `json:"user"`
	Images     []ImagesDetail `json:"images"`
}

type UserDetail struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type ImagesDetail struct {
	ImageUrl  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

// singular
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campFormatter := CampaignFormatter{
		Id:         campaign.ID,
		UserId:     campaign.UserId,
		Name:       campaign.Name,
		ShortDesc:  campaign.ShortDescription,
		ImageUrl:   "",
		GoalAmount: campaign.GoalAmount,
		CurrAmount: campaign.CurrentAmount,
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

func FormatCampaignDetail(campaign Campaign) DetailCampaignFormatter {
	var detailCampaign = DetailCampaignFormatter{
		Id:         campaign.ID,
		UserId:     campaign.UserId,
		Name:       campaign.Name,
		ShortDesc:  campaign.ShortDescription,
		ImageUrl:   "",
		GoalAmount: campaign.GoalAmount,
		CurrAmount: campaign.CurrentAmount,
	}
	if len(campaign.CampaignImages) > 0 {
		detailCampaign.ImageUrl = campaign.CampaignImages[0].FileName
	}
	var perks []string
	perks = append(perks, strings.Split(campaign.Perks, ",")...)
	detailCampaign.Perks = perks

	// add user
	var userDetail = UserDetail{
		campaign.User.Name,
		campaign.User.AvatarFileName,
	}
	detailCampaign.User = userDetail

	// images
	images := []ImagesDetail{}

	for _, image := range campaign.CampaignImages {
		campaignImageFormater := ImagesDetail{}
		campaignImageFormater.ImageUrl = image.FileName
		campaignImageFormater.IsPrimary = false
		if image.IsPrimary == 1 {
			campaignImageFormater.IsPrimary = true
		}
		images = append(images, campaignImageFormater)
	}
	detailCampaign.Images = images
	return detailCampaign
}
