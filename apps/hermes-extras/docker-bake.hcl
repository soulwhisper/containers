DATE   = formatdate("YYYY.MM.DD", timestamp())
APP    = "hermes-extras"
SOURCE = "https://github.com/soulwhisper/containers"
variable "GIT_SHA" {}

# Feishu platform deps for the hermes-agent runtime (lark-oapi + its one dep
# missing from the upstream venv, requests-toolbelt). Mirrored onto the data
# PVC by install.sh and loaded via PYTHONPATH — runtime lazy-install is
# impossible: read-only rootfs, and plugin deps belong in GitOps content.
variable "LARK_OAPI_VERSION" {
  // renovate: datasource=pypi depName=lark-oapi
  default = "1.6.8"
}

variable "REQUESTS_TOOLBELT_VERSION" {
  // renovate: datasource=pypi depName=requests-toolbelt
  default = "1.0.0"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    LARK_OAPI_VERSION         = "${LARK_OAPI_VERSION}"
    REQUESTS_TOOLBELT_VERSION = "${REQUESTS_TOOLBELT_VERSION}"
  }
  labels = {
    "org.opencontainers.image.vendor"   = "soulwhisper"
    "org.opencontainers.image.source"   = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created"  = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title"    = "${APP}"
    "org.opencontainers.image.url"      = "${SOURCE}"
    "org.opencontainers.image.version"  = "${DATE}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output   = ["type=docker"]
  tags     = ["${APP}:${DATE}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${DATE}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
