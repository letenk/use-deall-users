package service

import (
	"errors"
	"fmt"

	"github.com/letenk/use_deal_user/models/domain"
	"github.com/letenk/use_deal_user/models/web"
	"github.com/letenk/use_deal_user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAll() ([]domain.User, error)
	GetOne(id string) (domain.User, error)
	Create(req web.UserCreateRequest) (string, error)
	Update(id string, req web.UserUpdateRequest) (bool, error)
	Delete(id string) (bool, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewServiceUser(repository repository.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) GetAll() ([]domain.User, error) {

	// Get all
	users, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return users, err
}

func (s *userService) GetOne(id string) (domain.User, error) {
	// Get one
	user, err := s.repository.GetOne(id)
	if err != nil {
		return user, err
	}

	return user, err
}

func (s *userService) Create(req web.UserCreateRequest) (string, error) {
	// Create new user object
	newUser := domain.User{}
	newUser.Fullname = req.Fullname
	newUser.Username = req.Username

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	newUser.Password = string(passwordHash)

	// If field role is empty string
	if req.Role == "" {
		// Ya, Set role as user
		newUser.Role = "user"
	} else {
		// No, take data field role
		newUser.Role = req.Role
	}

	// Insert
	userId, err := s.repository.Insert(newUser)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *userService) Update(id string, req web.UserUpdateRequest) (bool, error) {
	// Get One
	user, err := s.repository.GetOne(id)
	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with ID %s Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	// Change fullname when field not empty
	if req.Fullname != "" {
		user.Fullname = req.Fullname
	}

	// Change password when field not empty
	if req.Password != "" {
		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
		if err != nil {
			return false, err
		}
		user.Password = string(passwordHash)
	}

	// Change role when field not empty
	if req.Role != "" {
		user.Role = req.Role
	}

	// Update
	ok, err := s.repository.Update(user)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *userService) Delete(id string) (bool, error) {
	// Get One
	user, err := s.repository.GetOne(id)
	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with ID %s Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	// Delete
	ok, err := s.repository.Delete(user.ID)
	if err != nil {
		return false, err
	}

	return ok, nil

}
