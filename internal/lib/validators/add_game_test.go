package validators

import (
	"testing"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/date"
)

func TestAddGame_validation(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		request         *game.AddGameRequest
		expectedValid   bool
		expectedMessage string
	}{
		{
			name: "no title",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: &date.Date{Year: 2018, Month: 12, Day: 25},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.TitleRequiredMessage,
		},
		{
			name: "no description",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "",
				ReleaseDate: &date.Date{Year: 2018, Month: 12, Day: 25},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.DescriptionRequiredMessage,
		},
		{
			name: "no date",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: nil,
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.ReleaseYearRequiredMessage,
		},
		{
			name: "no year",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: &date.Date{Year: 2019},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.ReleaseYearRequiredMessage,
		},
		{
			name: "no month",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: &date.Date{Year: 2019, Day: 2},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.ReleaseYearRequiredMessage,
		},
		{
			name: "no day",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: &date.Date{Year: 2019, Month: 2},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   false,
			expectedMessage: outerror.ReleaseYearRequiredMessage,
		},
		{
			name: "valid",
			request: &game.AddGameRequest{Game: &game.GameRequest{
				Title:       "qwe",
				Genres:      nil,
				Description: "qwe",
				ReleaseDate: &date.Date{Year: 2018, Month: 12, Day: 25},
				CoverImage:  nil,
				Tags:        nil,
			}},
			expectedValid:   true,
			expectedMessage: "",
		},
	}
	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			t.Parallel()
			gotValid, gotMessage := AddGame(ts.request)
			assert.Equal(t, ts.expectedValid, gotValid)
			assert.Equal(t, ts.expectedMessage, gotMessage)
		})
	}
}
