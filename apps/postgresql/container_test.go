package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/postgresql:18.6")

	lib := "/usr/lib/postgresql/18/lib/"

	t.Run("cnpg compatibility: postgres user is uid 26", func(t *testing.T) {
		testhelpers.TestUserUID(t, ctx, image, "postgres", 26)
	})
	t.Run("pgroonga library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pgroonga.so", nil)
	})
	t.Run("pgvector library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vector.so", nil)
	})

	t.Run("pgaudit library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pgaudit.so", nil)
	})

	t.Run("pg_failover_slots library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pg_failover_slots.so", nil)
	})

	t.Run("postgis library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"postgis-3.so", nil)
	})

	t.Run("postgres binary runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "postgres", "--version")
	})
}
