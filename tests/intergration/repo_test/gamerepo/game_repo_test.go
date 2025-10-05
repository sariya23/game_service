//go:build integrations

package gamerepo

import (
	"context"
	"testing"
)

func TestSaveGame_success(t *testing.T) {
	ctx := context.Background()
	db.SetUp(ctx, t, tables...)
	defer db.TearDown(t)
}
