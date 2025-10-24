package dto

import (
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
	"slices"
	"strings"
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

type QueryMovie struct {
	Title   string
	Genres  []string
	Filters Filters
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f *Filters) Limit() int {
	return f.PageSize
}

func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f *Filters) SortColumn() string {
	if slices.Contains(f.SortSafeList, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f *Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitzero"`
	PageSize     int `json:"page_size,omitzero"`
	FirstPage    int `json:"first_page,omitzero"`
	LastPage     int `json:"last_page,omitzero"`
	TotalRecords int `json:"total_records,omitzero"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
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

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
