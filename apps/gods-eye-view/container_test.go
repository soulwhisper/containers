package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/gods-eye-view:latest")

	t.Run("Check runtime user is node", func(t *testing.T) {
		testhelpers.TestUserUID(t, ctx, image, "node", 1000)
	})

	// ---- upstream checkout ------------------------------------------------

	for _, f := range []string{
		"/app/vite.config.js",
		"/app/server/standalone/vite.config.js",
		"/app/server/providers/local.js",
		"/app/package-lock.json",
	} {
		f := f
		t.Run("Check file exists: "+f, func(t *testing.T) {
			testhelpers.TestFileExists(t, ctx, image, f, nil)
		})
	}

	// Puppeteer is a QA-only dep; its Chromium download is skipped at build.
	t.Run("Check puppeteer chromium not downloaded", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "! test -d /root/.cache/puppeteer")
	})

	// ---- runtime defaults -------------------------------------------------

	envDefaults := map[string]string{
		"HOST":                         "0.0.0.0",
		"PORT":                         "4173",
		"GEV_RATELIMIT_OPENAI_PER_MIN": "30",
		"GEV_RATELIMIT_GOOGLE_PER_MIN": "60",
	}
	for k, v := range envDefaults {
		key, val := k, v
		t.Run("Check env default: "+key, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil,
				"sh", "-c", `[ "$`+key+`" = "`+val+`" ]`)
		})
	}

	// ---- server boots and serves the app -----------------------------------

	// Without provider keys the app must still boot — it falls back to
	// keyless Esri imagery per upstream docs.
	t.Run("Check web UI serves", func(t *testing.T) {
		testhelpers.TestHTTPEndpoint(t, ctx, image, testhelpers.HTTPTestConfig{
			Port:       "4173",
			Path:       "/",
			StatusCode: 200,
		}, nil)
	})
}
