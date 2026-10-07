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

	// bridge.py constructs mqtt.Client(mqtt.CallbackAPIVersion.VERSION2) — a
	// paho 2.x-only API. A bare import passes on 1.x, so assert the contract.
	t.Run("paho v2 callback API contract", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"python", "-c", "from paho.mqtt.client import Client, CallbackAPIVersion; Client(CallbackAPIVersion.VERSION2)")
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

}
