package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/manyfold:latest")

	t.Run("health endpoint as non-root with read-only root filesystem", func(t *testing.T) {
		testhelpers.TestHTTPEndpoint(t, ctx, image, testhelpers.HTTPTestConfig{
			Port:       "3214",
			Path:       "/health",
			StatusCode: 200,
		}, &testhelpers.ContainerConfig{
			Env: map[string]string{
				"SECRET_KEY_BASE": "manyfold-container-test",
			},
			ReadOnlyRootfs: true,
			Tmpfs: map[string]string{
				"/config":          "rw,noexec,nosuid,size=128m,mode=0777",
				"/libraries/main":  "rw,noexec,nosuid,size=64m,mode=0777",
				"/tmp":             "rw,noexec,nosuid,size=128m,mode=0777",
				"/usr/src/app/tmp": "rw,noexec,nosuid,size=128m,mode=0777",
				"/usr/src/app/log": "rw,noexec,nosuid,size=64m,mode=0777",
			},
		})
	})
}
