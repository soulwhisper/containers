package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/manyfold:latest")

	paths := []string{
		"/usr/local/bin/manyfold-entrypoint.sh",
		"/usr/local/bin/manyfold-start.sh",
		"/usr/bin/supervisord",
		"/usr/bin/redis-server",
		"/etc/supervisord.conf",
	}

	for _, path := range paths {
		path := path
		t.Run("Check installed file: "+path, func(t *testing.T) {
			testhelpers.TestFileExists(t, ctx, image, path, nil)
		})
	}
}
