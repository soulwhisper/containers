package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/covel:latest")

	t.Run("default runtime user is node", func(t *testing.T) {
		testhelpers.TestUserUID(t, ctx, image, "node", 1000)
	})

	t.Run("health endpoint with read-only root filesystem", func(t *testing.T) {
		testhelpers.TestHTTPEndpoint(t, ctx, image, testhelpers.HTTPTestConfig{
			Port:       "3001",
			Path:       "/api/health",
			StatusCode: 200,
		}, &testhelpers.ContainerConfig{
			Env: map[string]string{
				"COVEL_BIND_HOST":          "0.0.0.0",
				"COVEL_HOME":               "/home/node/.covel",
				"COVEL_MEDIA_TOKEN_SECRET": "7e4f1ee2e05d28f6b52f66793a0f3f79843bb2b7a1f7b4e265f51eeeb45a73",
				"MEDIA_ROOT":               "/home/node/.covel/data/media",
				"SQLITE_PATH":              "/home/node/.covel/data/covel.db",
				"STORE_BACKEND":            "sqlite",
			},
			ReadOnlyRootfs: true,
			Tmpfs: map[string]string{
				"/home/node/.covel": "rw,noexec,nosuid,size=64m,mode=0777",
				"/tmp":              "rw,noexec,nosuid,size=16m,mode=1777",
			},
		})
	})
}
