package song

import "errors"

var (
	InvalidInputData    = errors.New("invalid input data")
	SongsNotFound       = errors.New("songs not found")
	InvalidSongIdFormat = errors.New("invalid song id format")
)
