package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/stretchr/testify/require"
)

func TestSaveMinio(t *testing.T) {
	ctx := context.Background()
	cfg := config.MustLoadByPath("../../../../config/local.env")
	s3 := MustPrepareMinio(ctx, mockslog.NewDiscardLogger(), cfg.Minio, false)
	t.Run("Успешное сохранение", func(t *testing.T) {
		image, err := random.Image()
		require.NoError(t, err)
		key, err := s3.SaveObject(ctx, fmt.Sprintf("test_%v", time.Now()), bytes.NewReader(image))
		require.NoError(t, err)
		require.NotEmpty(t, key)
	})
}
