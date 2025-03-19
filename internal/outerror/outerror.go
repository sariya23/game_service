package outerror

import "errors"

var (
	ErrGameAlreadyExist = errors.New("game with this title already exist")
)
