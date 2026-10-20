DATE = formatdate( "YYYY.MM.DD", timestamp() )
APP = "dsh-web"
SOURCE = "https://github.com/deepseek-ai/deepseek-harness"
variable "GIT_SHA" {}

# Image tags track the upstream dsh version (the workbench's upstream).
# Same depName+datasource as the base Dockerfile ARG and base bake var →
# renovate ships all of them in one grouped PR.
variable "DSH_VERSION" {
  // renovate: datasource=npm depName=@deepseek-ai/dsh
  default = "0.2.1-rc.2"
}

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
  tags = ["${APP}:${DSH_VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64"
  ]
  tags = [
    "ghcr.io/soulwhisper/${APP}:sha-${GIT_SHA}",
    "ghcr.io/soulwhisper/${APP}:${DSH_VERSION}",
    "ghcr.io/soulwhisper/${APP}:latest",
  ]
}

target "docker-metadata-action" {}
