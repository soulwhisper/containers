DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "dsh-tavern"
SOURCE = "https://github.com/flizzywine/dsh-tavern"
variable "GIT_SHA" {}

group "default" {
  targets = ["image-local"]
}

# Builds the shared basement first so `FROM deepseek-harness` resolves
# locally (no third-party registry dependency in CI or locally).
target "dsh-base" {
  context = "../deepseek-harness"
  tags = ["deepseek-harness:local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  contexts = {
    deepseek-harness = "target:dsh-base"
  }
  labels = {
    "org.opencontainers.image.vendor" = "soulwhisper"
    "org.opencontainers.image.source" = "https://github.com/soulwhisper/containers"
    "org.opencontainers.image.created" = "${DATE}"
    "org.opencontainers.image.revision" = "${GIT_SHA}"
    "org.opencontainers.image.title" = "${APP}"
    "org.opencontainers.image.url" = "${SOURCE}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:local"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
