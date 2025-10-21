package service

import (
	"context"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, error)
	GetMovieById(ctx context.Context, id int64) (*domain.Movie, error)
	UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type movieService struct {
	movieRepository repository.MovieRepository
}

func (m *movieService) CreateMovie(ctx context.Context, input *dto.Movie) (*domain.Movie, error) {
	var movie *domain.Movie
	if err := m.movieRepository.CreateMovie(ctx, &domain.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}); err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *movieService) GetMovieById(ctx context.Context, id int64) (*domain.Movie, error) {
	return m.movieRepository.GetMovieById(ctx, id)
}

func (m *movieService) UpdateMovie(ctx context.Context, id int64, input *dto.UpdateMovie) (*domain.Movie, error) {
	movie, err := m.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	updatedMovie, err := m.movieRepository.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return updatedMovie, nil
}

func (m *movieService) DeleteMovie(ctx context.Context, id int64) error {
	return m.movieRepository.DeleteMovie(ctx, id)
}

func NewMovieService(movieRepository repository.MovieRepository) MovieService {
	return &movieService{
		movieRepository: movieRepository,
	}
}
