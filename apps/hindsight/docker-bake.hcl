DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "hindsight"
SOURCE = "https://github.com/vectorize-io/hindsight"
variable "GIT_SHA" {}

variable "VERSION" {
  // renovate: datasource=docker depName=ghcr.io/vectorize-io/hindsight versioning=docker
  default = "0.8.4-slim"
}

group "default" {
  targets = ["image-local"]
}

variable "VERSION_TRIM" {
  default = trimsuffix(trimprefix(VERSION, "v"), "-slim")
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
    "org.opencontainers.image.version" = "${VERSION_TRIM}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION_TRIM}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${VERSION_TRIM}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]

}

target "docker-metadata-action" {}
