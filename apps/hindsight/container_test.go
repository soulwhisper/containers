package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hindsight:latest")

	python := "/app/api/.venv/bin/python"

	// ---- venv & local-ml deps ---------------------------------------------

	t.Run("Check hindsight api venv python exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, python, nil)
	})

	t.Run("Check torch is the CPU build", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, python, "-c",
			"import torch; assert torch.version.cuda is None, torch.version.cuda")
	})

	for _, mod := range []string{"sentence_transformers", "transformers"} {
		m := mod
		t.Run("Check module importable: "+m, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil, python, "-c", "import "+m)
		})
	}

	// ---- baked-in model weights -------------------------------------------

	t.Run("Check HF cache dir exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-d", "/app/models/hub")
	})

	for _, snap := range []string{
		"models--BAAI--bge-m3",
		"models--BAAI--bge-reranker-v2-m3",
	} {
		s := snap
		t.Run("Check model snapshot baked in: "+s, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil,
				"test", "-d", "/app/models/hub/"+s)
		})
	}

	// ---- runtime local-provider defaults ----------------------------------

	envDefaults := map[string]string{
		"HF_HOME":                              "/app/models",
		"HINDSIGHT_API_EMBEDDINGS_PROVIDER":    "local",
		"HINDSIGHT_API_EMBEDDINGS_LOCAL_MODEL": "BAAI/bge-m3",
		"HINDSIGHT_API_RERANKER_PROVIDER":      "local",
		"HINDSIGHT_API_RERANKER_LOCAL_MODEL":   "BAAI/bge-reranker-v2-m3",
	}
	for k, v := range envDefaults {
		key, val := k, v
		t.Run("Check env default: "+key, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil,
				"sh", "-c", `[ "$`+key+`" = "`+val+`" ]`)
		})
	}
}
