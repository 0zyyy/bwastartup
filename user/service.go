package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUserInput) (User, error)
	CheckEmailAvail(email EmailUserInput) (bool, error)
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
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "User"
	newUser, err := s.repository.AddUser(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) Login(input LoginUserInput) (User, error) {
	user, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("no user found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return user, errors.New("password salah")
	}
	return user, nil
}

func (s *service) CheckEmailAvail(email EmailUserInput) (bool, error) {
	user, err := s.repository.FindByEmail(email.Email)
	if err != nil {
		return false, err
	}
	if user.Id == 0 {
		return true, nil
	}
	return false, nil
}
