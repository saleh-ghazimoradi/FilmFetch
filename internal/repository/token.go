package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"time"
)

type TokenRepository interface {
	InsertToken(ctx context.Context, token *domain.Token) error
	DeleteAllForUser(ctx context.Context, scope string, userId int64) error
	GetForToken(ctx context.Context, scope, plainText string) (*domain.User, error)
}

type tokenRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
}

func (t *tokenRepository) InsertToken(ctx context.Context, token *domain.Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := t.dbWrite.ExecContext(ctx, query, args...)
	return err
}

func (t *tokenRepository) DeleteAllForUser(ctx context.Context, scope string, userId int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := t.dbWrite.ExecContext(ctx, query, scope, userId)
	return err
}

func (t *tokenRepository) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user domain.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := t.dbRead.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func NewTokenRepository(dbWrite, dbRead *sql.DB) TokenRepository {
	return &tokenRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
