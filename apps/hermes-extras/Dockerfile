# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# buzz CLI — required by the hermes `buzz` platform plugin (adapter shells out
# to `buzz`, JSON in/out). Upstream ships no CLI release assets and the relay
# image (ghcr.io/block/buzz) carries only buzz-relay/buzz-admin/buzz-pair-relay,
# so build buzz-cli from source at the pinned block/buzz git tag (BUZZ_VERSION,
# renovate-tracked via github-tags).
FROM rust:1.88-trixie AS buzz-cli
ARG BUZZ_VERSION
RUN set -eux; \
    git clone --depth 1 --branch "${BUZZ_VERSION}" https://github.com/block/buzz /build; \
    cd /build; \
    cargo build --release --locked -p buzz-cli; \
    strip target/release/buzz

# https://github.com/NousResearch/hermes-agent/blob/main/Dockerfile
FROM debian:13-slim@sha256:3a39a0592364683e6bab97937b72cad5a8fa6dcbbee90edb3bb48c7f8e94f258

RUN apt update -y && \
    apt install -y --no-install-recommends curl ca-certificates && \
    install -dm 755 /etc/apt/keyrings && \
    curl -fSs https://mise.en.dev/gpg-key.pub | tee /etc/apt/keyrings/mise-archive-keyring.asc 1> /dev/null  && \
    echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.asc] https://mise.en.dev/deb stable main" | tee /etc/apt/sources.list.d/mise.list  && \
    apt update -y && apt install -y --no-install-recommends mise file && \
    rm -rf /var/lib/apt/lists/*

# install every tool declared in .mise.toml
# github_token secret: MISE_GITHUB_TOKEN authenticates api.github.com release
# resolution (anonymous = 60 req/h); optional, local builds fall back to anonymous
COPY app /app
WORKDIR /app
RUN --mount=type=secret,id=github_token,required=false \
    MISE_GITHUB_TOKEN="$(cat /run/secrets/github_token 2>/dev/null || true)"; \
    export MISE_GITHUB_TOKEN; \
    mise trust && mise install

# Collect binaries into a flat /data/ layout.
# mise bin-paths prints one path per installed tool (handles bin/ vs root layouts).
RUN set -eux; \
    mkdir -p /data; \
    for bin_dir in $(mise bin-paths); do \
      find "$bin_dir" -maxdepth 1 -type f -perm -u+x \
        ! -name '*.txt' ! -name '*.md' ! -name '*.sh' \
        -exec sh -c 'file "$1" | grep -q "ELF.*executable" && cp -v "$1" /data/' _ {} \; ; \
    done; \
    chown -R 10000:10000 /data

# buzz-cli is built in its own stage (see top of file); add it to the bundle.
COPY --chown=10000:10000 --from=buzz-cli /build/target/release/buzz /data/buzz

# caveman skill pack — prompt skills, not binaries. Upstream ships no release
# assets, so mise can't fetch it; pull the pinned release tarball instead.
ARG CAVEMAN_VERSION
RUN set -eux; \
    curl -fsSL "https://github.com/JuliusBrussee/caveman/archive/refs/tags/${CAVEMAN_VERSION}.tar.gz" -o /tmp/caveman.tar.gz; \
    tar -xzf /tmp/caveman.tar.gz -C /tmp; \
    mkdir -p /skills; \
    cp -r "/tmp/caveman-${CAVEMAN_VERSION#v}/skills/." /skills/; \
    rm -rf /tmp/caveman*; \
    chown -R 10000:10000 /skills

USER 10000:10000
WORKDIR /data

# No ENTRYPOINT: init-containers and tests override `command` fully.
# Default CMD prints version table — useful for `docker run` discovery.
CMD ["/bin/sh", "-c", "for f in /data/*; do printf '=== %s ===\\n' \"$(basename \"$f\")\"; \"$f\" --version 2>&1 || true; done"]
