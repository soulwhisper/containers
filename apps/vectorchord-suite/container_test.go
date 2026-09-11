package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/vectorchord-suite:latest")

	lib := "/usr/lib/postgresql/18/lib/"
	ext := "/usr/share/postgresql/18/extension/"

	t.Run("vchord library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vchord.so", nil)
	})

	t.Run("vchord extension sql exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "ls "+ext+"vchord--*.sql >/dev/null")
	})

	t.Run("pgroonga library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"pgroonga.so", nil)
	})

	t.Run("pgroonga extension sql exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "ls "+ext+"pgroonga--*.sql >/dev/null")
	})

	t.Run("vchord_bm25 library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vchord_bm25.so", nil)
	})

	t.Run("pg_tokenizer library exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "ls "+lib+"pg_tokenizer.so >/dev/null || ls "+lib+"vchord_pg_tokenizer.so >/dev/null")
	})

	t.Run("pgvector library exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, lib+"vector.so", nil)
	})

	t.Run("postgres binary runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "postgres", "--version")
	})
}
