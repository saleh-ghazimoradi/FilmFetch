package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"time"
)

func GenerateToken(userId int64, ttl time.Duration, scope string) *domain.Token {
	token := &domain.Token{
		Plaintext: rand.Text(),
		UserId:    userId,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token
}
