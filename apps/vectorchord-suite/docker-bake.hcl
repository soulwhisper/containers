DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "vectorchord-suite"
SOURCE = "https://github.com/soulwhisper/containers"
variable "GIT_SHA" {}

# pig CLI version (github:pgsty/pig release). Extension versions are
# resolved by pig per PG major — nothing else to pin here.
variable "PIG_VERSION" {
  // renovate: datasource=github-releases depName=pgsty/pig
  default = "1.8.1"
}

variable "TAG_MAIN" {
  default = "pig${PIG_VERSION}-pg18"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    PIG_VERSION = "${PIG_VERSION}"
  }
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
    "org.opencontainers.image.version" = "${TAG_MAIN}"
    "vectorchord-suite.pig-version" = "${PIG_VERSION}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${TAG_MAIN}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${TAG_MAIN}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
