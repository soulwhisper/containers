DATE   = formatdate("YYYY.MM.DD", timestamp())
APP    = "covel"
SOURCE = "https://github.com/ackness/covel"
variable "GIT_SHA" {}

variable "VERSION" {
  // renovate: datasource=github-releases depName=ackness/covel
  default = "v0.0.54"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits  = ["docker-metadata-action"]
  context   = "https://github.com/ackness/covel.git?tag=${VERSION}"
  dockerfile = "docker/Dockerfile"
  labels = {
    "org.opencontainers.image.vendor"  = "soulwhisper"
    "org.opencontainers.image.source"  = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title"   = "${APP}"
    "org.opencontainers.image.url"     = "${SOURCE}"
    "org.opencontainers.image.version" = "${trimprefix(VERSION, "v")}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output   = ["type=docker"]
  tags     = ["${APP}:${trimprefix(VERSION, "v")}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${trimprefix(VERSION, "v")}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
