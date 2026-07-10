DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "dify-web-client"
SOURCE = "https://github.com/langgenius/webapp-conversation"
variable "GIT_SHA" {}

variable "VERSION" {
  // renovate: datasource=git-refs depName=https://github.com/langgenius/webapp-conversation
  default = "33085b6608fe7174e0fb75e46c220348863b1c19"
}

# Short hash from the full commit SHA for readable tags.
variable "VERSION_SHORT" {
  default = substr(VERSION, 0, 8)
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
  tags = ["${APP}:${VERSION_SHORT}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${VERSION_SHORT}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]

}

target "docker-metadata-action" {}
