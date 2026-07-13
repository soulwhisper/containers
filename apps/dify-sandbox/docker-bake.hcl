DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "dify-sandbox"
SOURCE = "https://github.com/langgenius/dify-sandbox"
variable "GIT_SHA" {}

# Upstream dify-sandbox base image (semver tag).
variable "VERSION" {
  // renovate: datasource=docker depName=langgenius/dify-sandbox versioning=semver
  default = "0.2.15"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    VERSION = "${VERSION}"
  }
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
    "org.opencontainers.image.version" = "${VERSION}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${VERSION}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
