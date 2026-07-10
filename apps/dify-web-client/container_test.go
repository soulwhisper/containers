package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dify-web-client:latest")

	t.Run("Check node exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/local/bin/node", nil)
	})

	t.Run("Check server.js exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/server.js", nil)
	})

	t.Run("Check static dir exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"test", "-d", "/app/.next/static")
	})

	t.Run("Check public dir exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"test", "-d", "/app/public")
	})

	t.Run("Check server.js is valid JS", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"node", "-c", "/app/server.js")
	})

}
