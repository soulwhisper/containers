package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/cloudnative-custom:18.6")

	lib := "/usr/lib/postgresql/18/lib/"
	ext := "/usr/share/postgresql/18/extension/"

	t.Run("vchord library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vchord.so", nil)
	})

	t.Run("vchord_bm25 library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vchord_bm25.so", nil)
	})

	t.Run("pg_tokenizer library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pg_tokenizer.so", nil)
	})

  t.Run("pgvector is absent (auto-detect consumers must use vchord)", func(t *testing.T) {
    testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "test ! -f "+lib+"vector.so")
  })

	t.Run("pgroonga library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pgroonga.so", nil)
	})
	t.Run("pgvector library exists (required by vchord)", func(t *testing.T) {
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

	t.Run("vchord extension sql exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "ls "+ext+"vchord--*.sql >/dev/null")
	})

	t.Run("postgres binary runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "postgres", "--version")
	})
}
