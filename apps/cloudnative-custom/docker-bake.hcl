DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "cloudnative-custom"
SOURCE = "https://github.com/soulwhisper/containers"
variable "GIT_SHA" {}

variable "PG_VERSION" {
  // renovate: datasource=docker depName=postgres versioning=semver
  default = "18.6"
}

# Major floater tag is derived from PG_VERSION
variable "PG_MAJOR" {
  default = split(".", PG_VERSION)[0]
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
    "org.opencontainers.image.version" = "${PG_VERSION}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${PG_VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${PG_VERSION}",
    "ghcr.io/soulwhisper/${APP}:${PG_MAJOR}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
