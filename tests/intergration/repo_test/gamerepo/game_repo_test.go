//go:build integrations

package gamerepo

import (
	"context"
	"testing"
)

func TestSaveGame(t *testing.T) {
	ctx := context.Background()
	db.SetUp(ctx, t, tables...)
	defer db.TearDown(t)
}
