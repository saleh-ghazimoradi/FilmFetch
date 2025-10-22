package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"time"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) error
	GetMovieById(ctx context.Context, id int64) (*domain.Movie, error)
	UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type movieRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
}

func (m *movieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	query := `INSERT INTO movies(title, year, runtime, genres) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.dbWrite.QueryRowContext(ctx, query, args...).Scan(&movie.Id, &movie.CreatedAt, &movie.Version)
}

func (m *movieRepository) GetMovieById(ctx context.Context, id int64) (*domain.Movie, error) {
	var movie domain.Movie
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, year, runtime, genres, version FROM movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.dbRead.QueryRowContext(ctx, query, id).Scan(&movie.Id, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m *movieRepository) UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	query := `
        UPDATE movies 
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.Id, movie.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.dbWrite.QueryRowContext(ctx, query, args...).Scan(&movie.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEditConflict
		default:
			return nil, err
		}
	}

	return movie, nil
}

func (m *movieRepository) DeleteMovie(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.dbWrite.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func NewMovieRepository(dbWrite, dbRead *sql.DB) MovieRepository {
	return &movieRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
