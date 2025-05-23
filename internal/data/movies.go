package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"whitelight.quigley.net/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (mm *MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO Movies(title, year, runtime, genres)
			 VALUES($1, $2, $3, $4)
			RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array([]int{123, 12})}

	return mm.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (mm *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM Movies WHERE id = $1`
	result, e := mm.DB.Exec(query, id)
	if e != nil {
		return e
	}
	rowsAffected, e := result.RowsAffected()
	if e != nil {
		return e
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (mm *MovieModel) Update(movie *Movie) error {
	query := `UPDATE Movies 
			  SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			  WHERE id = $5
			  RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return mm.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (mm *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, title, year, runtime, genres, version
			  FROM Movies WHERE id = $1`
	var movie Movie

	e := mm.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if e != nil {
		switch {
		case errors.Is(e, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, e
		}
	}
	return &movie, nil
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
