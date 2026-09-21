DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "gods-eye-view"
SOURCE = "https://github.com/bilawalsidhu/gods-eye-view"
variable "GIT_SHA" {}

# Upstream git ref. Upstream has no release tags yet and moves fast on main;
# track the branch head (floating, like a latest tag) until tags exist.
variable "VERSION" {
  default = "main"
}

group "default" {
  targets = ["image-local"]
}

variable "VERSION_TRIM" {
  default = trimprefix(VERSION, "v")
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
