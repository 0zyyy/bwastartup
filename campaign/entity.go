package campaign

import "time"

type Campaign struct {
	ID             int
	UserId         int
	Name           string
	ShortDesc      string
	Desc           string
	Perks          string
	BackerCount    int
	GoalAmount     int
	CurrAmount     int
	Slug           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CampaignImages []CampaignImage
}

type CampaignImage struct {
	Id         int
	CampaignId int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
