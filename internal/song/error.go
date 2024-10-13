package song

import "errors"

var (
	InvalidInputData      = errors.New("invalid input data")
	SongsNotFound         = errors.New("songs not found")
	AuthorNotFound        = errors.New("author not found")
	AuthorAlreadyExists   = errors.New("author already exists")
	AuthorSongDuplicate   = errors.New("song with this name already exists by this author")
	InvalidSongIdFormat   = errors.New("invalid song id format")
	InvalidAuthorIdFormat = errors.New("invalid author id format")
	ErrorGetSongData      = errors.New("error get song data")
	ErrorGetSongLyrics    = errors.New("error get song lyrics")
)
