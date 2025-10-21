package dto

import (
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
	"time"
)

type Movie struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

type UpdateMovie struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *int32   `json:"runtime"`
	Genres  []string `json:"genres"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

func ValidateUpdateMovie(v *validator.Validator, update *UpdateMovie) {
	if update.Title != nil {
		v.Check(*update.Title != "", "title", "must be provided")
		v.Check(len(*update.Title) <= 500, "title", "must not be more than 500 bytes long")
	}

	if update.Year != nil {
		v.Check(*update.Year != 0, "year", "must be provided")
		v.Check(*update.Year >= 1888, "year", "must be greater than 1888")
		v.Check(*update.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	}

	if update.Runtime != nil {
		v.Check(*update.Runtime != 0, "runtime", "must be provided")
		v.Check(*update.Runtime > 0, "runtime", "must be a positive integer")
	}

	if update.Genres != nil {
		v.Check(len(update.Genres) >= 1, "genres", "must contain at least 1 genre")
		v.Check(len(update.Genres) <= 5, "genres", "must not contain more than 5 genres")
		v.Check(validator.Unique(update.Genres), "genres", "must not contain duplicate values")
	}
}
