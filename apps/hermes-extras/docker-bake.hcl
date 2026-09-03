DATE   = formatdate("YYYY.MM.DD", timestamp())
APP    = "hermes-extras"
SOURCE = "https://github.com/soulwhisper/containers"
variable "GIT_SHA" {}

# caveman skill pack release tag (no release assets upstream -> tarball).
variable "CAVEMAN_VERSION" {
  // renovate: datasource=github-releases depName=JuliusBrussee/caveman
  default = "v2.5.0"
}

# buzz relay image tag the bundled CLI is paired with (renovate tracks the tag;
# BUZZ_REF below must be updated to that tag's org.opencontainers.image.revision).
variable "BUZZ_VERSION" {
  // renovate: datasource=docker depName=ghcr.io/block/buzz
  default = "0.1.1"
}

# source revision of ghcr.io/block/buzz:${BUZZ_VERSION}
variable "BUZZ_REF" {
  default = "68a0cc8506be4ea1fb110b65eec0787e4cd84378"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    CAVEMAN_VERSION = "${CAVEMAN_VERSION}"
    BUZZ_REF        = "${BUZZ_REF}"
  }
  # read from the bake-action step env (MISE_GITHUB_TOKEN)
  secret = ["id=github_token,env=MISE_GITHUB_TOKEN"]
  labels = {
    "org.opencontainers.image.vendor"   = "soulwhisper"
    "org.opencontainers.image.source"   = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created"  = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title"    = "${APP}"
    "org.opencontainers.image.url"      = "${SOURCE}"
    "org.opencontainers.image.version"  = "${DATE}"
    "hermes-extras.buzz-version"        = "${BUZZ_VERSION}"
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
