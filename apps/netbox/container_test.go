package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/netbox-custom:latest")

	pip := "/opt/netbox/venv/bin/pip"

	t.Run("Check netbox venv python exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/netbox/venv/bin/python", nil)
	})

	t.Run("Check netbox venv pip exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, pip, nil)
	})

	plugins := []string{
		"netbox-attachments",
		"netbox-bgp",
		"netbox-floorplan-plugin",
		"netbox-interface-synchronization",
		"netbox-plugin-dns",
		"netbox-plugin-prometheus-sd",
		"netbox-qrcode",
		"netbox-reorder-rack",
		"netbox-topology-views",
		"netbox-napalm-plugin",
	}

	for _, p := range plugins {
		pkg := p
		t.Run("Check plugin installed: "+pkg, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil, pip, "show", pkg)
		})
	}
}
