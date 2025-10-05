package validators

import (
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/date"
)

func TestAddGame_validation(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		request         *gamev2.AddGameRequest
		expectedValid   bool
		expectedMessage string
	}{
		{
			name: "no title",
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
			request: &gamev2.AddGameRequest{Game: &gamev2.GameRequest{
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
