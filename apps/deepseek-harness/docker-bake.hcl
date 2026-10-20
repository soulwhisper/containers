DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "deepseek-harness"
SOURCE = "https://github.com/deepseek-ai/deepseek-harness"
variable "GIT_SHA" {}

# Upstream DeepSeek Harness version (npm @deepseek-ai/dsh). Plugin versions
# and the tavern ref are pinned as Dockerfile ARGs with renovate comments;
# a renovate group ships related bumps in one PR where possible.
# renovate: datasource=npm depName=@deepseek-ai/dsh
variable "DSH_VERSION" {
  default = "0.2.1-rc.2"
}

group "default" {
  targets = ["image-local"]
}

variable "DSH_VERSION_TRIM" {
  default = trimprefix(DSH_VERSION, "v")
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
    "org.opencontainers.image.version" = "${DSH_VERSION_TRIM}"
  }
  no-cache = true
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${DSH_VERSION_TRIM}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${DSH_VERSION_TRIM}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
