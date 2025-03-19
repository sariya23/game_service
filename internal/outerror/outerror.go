package outerror

import "errors"

var (
	ErrGameAlreadyExist = errors.New("game with this title already exist")
)

var (
	TitleRequiredMessage       = "Title is required field"
	DescriptionRequiredMessage = "Description is required field"
	ReleaseYearRequiredMessage = "Release Year is required field"
)
