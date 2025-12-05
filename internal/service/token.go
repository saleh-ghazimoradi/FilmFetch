package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
	"github.com/saleh-ghazimoradi/FilmFetch/utils"
	"time"
)

type TokenService interface {
	Tokenize(ctx context.Context, userId int64, ttl time.Duration, scope string) (*domain.Token, error)
	ActivateUser(ctx context.Context, input *dto.ActivateUser) (*domain.User, error)
}

type tokenService struct {
	tokenRepository repository.TokenRepository
	userRepository  repository.UserRepository
}

func (t *tokenService) Tokenize(ctx context.Context, userId int64, ttl time.Duration, scope string) (*domain.Token, error) {
	token := utils.GenerateToken(userId, ttl, scope)

	err := t.tokenRepository.InsertToken(ctx, token)
	return token, err
}

func (t *tokenService) ActivateUser(ctx context.Context, input *dto.ActivateUser) (*domain.User, error) {
	user, err := t.tokenRepository.GetForToken(ctx, domain.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		return nil, err
	}

	user.Activated = true

	if err := t.userRepository.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEditConflict) {
			return nil, repository.ErrEditConflict
		}
		return nil, err
	}

	if err := t.tokenRepository.DeleteAllForUser(ctx, domain.ScopeActivation, user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func NewTokenService(tokenRepository repository.TokenRepository, userRepository repository.UserRepository) TokenService {
	return &tokenService{
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
	}
}
