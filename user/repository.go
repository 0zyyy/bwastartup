package user

import "gorm.io/gorm"

type Repository interface {
	AddUser(user User) (User, error)
	FindByEmail(email string) (User, error)
	UploadAvatar(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddUser(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (r *repository) UploadAvatar(user User) (User, error) {
	avatarName := user.AvatarFileName
	r.db.Create(avatarName)
	return user, nil
}
