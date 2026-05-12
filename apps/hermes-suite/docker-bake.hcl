DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "hermes-suite"
SOURCE = "https://github.com/NousResearch/hermes-agent"
variable "GIT_SHA" {}

# Upstream hermes-agent base image (date-based tag).
variable "AGENT_VERSION" {
  // renovate: datasource=docker depName=nousresearch/hermes-agent
  default = "v2026.5.7"
}

# Upstream hermes-webui release tag (semver).
variable "WEBUI_VERSION" {
  // renovate: datasource=github-releases depName=nesquena/hermes-webui
  default = "v0.51.44"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    AGENT_VERSION = "${AGENT_VERSION}"
    WEBUI_VERSION = "${WEBUI_VERSION}"
  }
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
    "org.opencontainers.image.version" = "${trimprefix(AGENT_VERSION, "v")}-${trimprefix(WEBUI_VERSION, "v")}"
    "hermes-suite.agent-version" = "${AGENT_VERSION}"
    "hermes-suite.webui-version" = "${WEBUI_VERSION}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${trimprefix(AGENT_VERSION, "v")}-${trimprefix(WEBUI_VERSION, "v")}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${trimprefix(AGENT_VERSION, "v")}-${trimprefix(WEBUI_VERSION, "v")}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]

}

target "docker-metadata-action" {}
