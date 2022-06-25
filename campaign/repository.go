package campaign

import "gorm.io/gorm"

type Repository interface {
	AddCampaign(campaign Campaign) (Campaign, error)
	GetCampaigns() ([]Campaign, error)
	FindByUserId(userId int) ([]Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) GetCampaigns() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserId(userId int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userId).Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
