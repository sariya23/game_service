package email

import (
	"testing"

	"github.com/sariya23/game_service/internal/config"
)

func TestSendMessage(t *testing.T) {
	cfg := config.MustLoadByPath("./config/local.env")
	sender := NewDialer()
}
