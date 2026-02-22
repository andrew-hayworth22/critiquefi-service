package types

import (
	"time"
)

type Media struct {
	ID                 int64
	MediaType          string
	Title              string
	ReleaseDate        time.Time
	Year               int
	Description        string
	ExternalReferences string
	CreatedAt          time.Time
	CreatedBy          int64
	UpdatedAt          NullableTime
	UpdatedBy          int64
}

type Film struct {
	Media
	FilmType       string
	RuntimeMinutes int
}

type Book struct {
	Media
	BookType string
	Pages    int
}

type Game struct {
	Media
	GameType string
}

type Music struct {
	Media
	MusicType     string
	LengthMinutes int
}

type Show struct {
	Media
	ShowType string
}
