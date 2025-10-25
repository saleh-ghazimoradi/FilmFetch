package service

import (
	"context"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, input *dto.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *dto.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) CreateUser(ctx context.Context, input *dto.User) (*domain.User, error) {
	us := domain.User{}
	us.Name = input.Name
	us.Email = input.Email
	if err := us.Password.Set(input.Password); err != nil {
		return nil, err
	}
	if err := u.userRepository.CreateUser(ctx, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepository.GetUserByEmail(ctx, email)
}

func (u *userService) UpdateUser(ctx context.Context, user *dto.User) error {
	return nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
