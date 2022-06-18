package user

import "golang.org/x/crypto/bcrypt"

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Nama
	user.Occupation = input.Occupation
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, nil
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "User"
	newUser, err := s.repository.AddUser(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}
