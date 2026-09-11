DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "postgresql"
SOURCE = "https://github.com/soulwhisper/containers"
variable "GIT_SHA" {}

# Release tags: postgres major (:18 floater) + major.minor.patch (:18.6 pin).
# cnpg clusters should pin the patch tag so pod restarts are reproducible;
# upstream pg updates flow through renovate PRs (version + digest).
# Major floater tag and the Dockerfile's extension install (-v ${PG_MAJOR})
# both derive from PG_VERSION (single renovate-tracked value).
variable "PG_MAJOR" {
  default = split(".", PG_VERSION)[0]
}

variable "PG_VERSION" {
  // renovate: datasource=docker depName=postgres versioning=semver
  default = "18.6"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    PG_MAJOR = "${PG_MAJOR}"
  }
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
