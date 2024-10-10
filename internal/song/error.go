package song

import "errors"

var (
	InvalidInputData    = errors.New("invalid input data")
	SongsNotFound       = errors.New("songs not found")
	InvalidSongIdFormat = errors.New("invalid song id format")
	ErrorGetSongIP      = errors.New("error get song ip")
	ErrorGetSongLyrics  = errors.New("error get song lyrics")
)
