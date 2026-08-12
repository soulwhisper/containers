# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

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
COPY app /app
WORKDIR /app
RUN mise trust && mise install

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
