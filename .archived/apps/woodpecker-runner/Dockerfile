# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# Woodpecker CI step runner for the home-ops repo.
#
# debian (glibc) base, deliberately not alpine:
# - python wheels without musllinux builds fall back to source builds on musl
#   (uvx/pip-based steps), glibc avoids the toolchain surprise
# - some release binaries ship glibc-only builds
# - docker CLI plugin lookup paths differ under busybox/musl setups
# Same apt-mise pattern as hermes-extras. Steps that need no tool installs
# keep their official images (docker:28-cli, astral-sh/uv, cnpg postgresql).
FROM debian:13-slim@sha256:3a39a0592364683e6bab97937b72cad5a8fa6dcbbee90edb3bb48c7f8e94f258

# git/curl: clone + forge API calls; jq: JSON munging;
# python3+pyyaml: label-sync/labeler scripts
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    ca-certificates curl git jq python3 python3-yaml file \
 && install -dm 755 /etc/apt/keyrings \
 && curl -fSs https://mise.en.dev/gpg-key.pub > /etc/apt/keyrings/mise-archive-keyring.asc \
 && echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.asc] https://mise.en.dev/deb stable main" > /etc/apt/sources.list.d/mise.list \
 && apt-get update \
 && apt-get install -y --no-install-recommends mise \
 && rm -rf /var/lib/apt/lists/*

# install every tool declared in app/.mise.toml
# github_token secret: MISE_GITHUB_TOKEN authenticates api.github.com release
# resolution (anonymous = 60 req/h); optional, local builds fall back to anonymous
COPY app /app
WORKDIR /app
RUN --mount=type=secret,id=github_token,required=false \
    MISE_GITHUB_TOKEN="$(cat /run/secrets/github_token 2>/dev/null || true)"; \
    export MISE_GITHUB_TOKEN; \
    mise trust && mise install

# Collect binaries onto PATH so steps run without mise shim overhead or
# runtime mise dependency.
RUN set -eux; \
    for bin_dir in $(mise bin-paths); do \
      find "$bin_dir" -maxdepth 1 -type f -perm -u+x \
        ! -name '*.txt' ! -name '*.md' \
        -exec sh -c 'file "$1" | grep -q "ELF.*executable" && cp -v "$1" /usr/local/bin/' _ {} \; ; \
    done

# root on purpose: woodpecker steps git-fetch/push and write kubeconfigs into
# the shared workspace volume; no USER directive.

# Default CMD prints the tool table — useful for `docker run` discovery.
CMD ["/bin/sh", "-c", "crane version; kubeconform -v; dnscontrol version; shellcheck --version | head -2; vendir version; talosctl version --client; kubectl version --client; yq --version; flate --version; jq --version; git --version; python3 --version"]
