package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/frigate-vision:latest")

	// ---- Bundled dependencies ----------------------------------------------

	t.Run("paho-mqtt imports", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"python", "-c", "import paho.mqtt.client")
	})

	t.Run("httpx imports", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"python", "-c", "import httpx")
	})

	// System site-packages, not user site: the cluster pod runs
	// readOnlyRootFilesystem + runAsUser 1000, so deps must be importable
	// without a writable HOME.
	t.Run("deps import with read-only HOME", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"HOME": "/nonexistent"}},
			"python", "-c", "import paho.mqtt.client, httpx")
	})

	// ---- Discovery -----------------------------------------------------------

	t.Run("Bundled versions are recorded", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"python", "-c", "import importlib.metadata as m; assert m.version('paho-mqtt'); assert m.version('httpx')")
	})
}
