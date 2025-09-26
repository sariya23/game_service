package outerror

import "errors"

var (
	ErrGameAlreadyExist           = errors.New("game with this title already exist")
	ErrGameNotFound               = errors.New("game not found")
	ErrTagNotFound                = errors.New("tag not found")
	ErrGenreNotFound              = errors.New("genre not found ")
	ErrCannotStartGameTransaction = errors.New("cannot start game transaction")
	ErrCannotSaveGameImage        = errors.New("cannot save image in s3")
	ErrImageNotFoundS3            = errors.New("game image not found in s3")
	ErrUnknownGameStatus          = errors.New("unknown game status")
	ErrInvalidNewGameStatus       = errors.New("invalid new game status")
)

var (
	TitleRequiredMessage              = "Title is required field"
	DescriptionRequiredMessage        = "Description is required field"
	ReleaseYearRequiredMessage        = "Release Year is required field"
	InternalMessage                   = "Internal error"
	GameAlreadyExistMessage           = "Game already exist"
	GameNotFoundMessage               = "Game not found"
	CannotStartGameTransactionMessage = "Cannot start game transaction"
	GameSavedWithoutImageMessage      = "Game saved but without image. Store is not response"
	GenreNotFoundMessage              = "Unknown genre name"
	TagNotFoundMessage                = "Unknown tag name"
	UnknownGameStatusMessage          = "Unknown game status"
	InvalidNewGameStatusMessage       = "Invalid new game status"
)
