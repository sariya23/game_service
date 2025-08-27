package postgresql

const (
	// Game table field names
	gameGameIDFieldName      = "game_id"
	gameTitleFieldName       = "title"
	gameDescriptionFieldName = "description"
	gameReleaseDateFieldName = "release_date"
	gameImageURLFieldName    = "image_url"

	// GameGenre table field names
	gameGenreGameIDFieldName  = "game_id"
	gameGenreGenreIDFieldName = "genre_id"

	// GameTag table field names
	gameTagGameIDFieldName = "game_id"
	gameTagTagIDFieldName  = "tag_id"

	// Genre table field names
	genreGenreIDFieldName   = "genre_id"
	genreGenreNameFieldName = "genre_name"

	// Tag table field names
	tagTagIDFieldName   = "tag_id"
	tagTagNameFieldName = "tag_name"
)
