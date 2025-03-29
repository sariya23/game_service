package outerror

import "errors"

var (
	ErrGameAlreadyExist           = errors.New("game with this title already exist")
	ErrGameNotFound               = errors.New("game not found")
	ErrCannotStartGameTransaction = errors.New("cannot start game transaction")
	ErrCannotSaveGameImage        = errors.New("cannot save image in s3")
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
)
