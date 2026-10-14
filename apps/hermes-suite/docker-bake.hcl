DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "hermes-suite"
SOURCE = "https://github.com/NousResearch/hermes-agent"
variable "GIT_SHA" {}

# Upstream hermes-agent release version. Since v0.21.6 upstream publishes
# semver GitHub releases instead of CalVer tags; Docker Hub carries only
# immutable rc.N-vX.Y.Z candidates plus the floating stable/latest aliases
# (no plain vX.Y.Z tag), so the version is tracked from GitHub releases here
# and the base image is digest-pinned separately in the Dockerfile.
variable "AGENT_VERSION" {
  // renovate: datasource=github-releases depName=NousResearch/hermes-agent
  default = "v0.21.7"
}

# Upstream hermes-webui release tag (semver).
variable "WEBUI_VERSION" {
  // renovate: datasource=github-releases depName=nesquena/hermes-webui
  default = "v0.52.113"
}

group "default" {
  targets = ["image-local"]
}

# local variables
variable "WEBUI_VERSION_TRIM" {
  default = trimprefix(WEBUI_VERSION, "v")
}
variable "AGENT_VERSION_TRIM" {
  default = trimprefix(AGENT_VERSION, "v")
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    WEBUI_VERSION = "${WEBUI_VERSION}"
  }
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
    "org.opencontainers.image.version" = "${AGENT_VERSION_TRIM}"
    "hermes-suite.agent-version" = "${AGENT_VERSION_TRIM}"
    "hermes-suite.webui-version" = "${WEBUI_VERSION_TRIM}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${AGENT_VERSION_TRIM}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${AGENT_VERSION_TRIM}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]

}

target "docker-metadata-action" {}
