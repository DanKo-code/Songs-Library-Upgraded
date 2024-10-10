package song

import "errors"

var (
	InvalidInputData    = errors.New("invalid input data")
	SongsNotFound       = errors.New("songs not found")
	SongAlreadyExists   = errors.New("song already exists")
	InvalidSongIdFormat = errors.New("invalid song id format")
	ErrorGetSongData    = errors.New("error get song data")
	ErrorGetSongLyrics  = errors.New("error get song lyrics")
)
