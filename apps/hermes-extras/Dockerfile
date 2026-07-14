# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

# https://github.com/NousResearch/hermes-agent/blob/main/Dockerfile
FROM debian:13-slim@sha256:020c0d20b9880058cbe785a9db107156c3c75c2ac944a6aa7ab59f2add76a7bd

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

USER 10000:10000
WORKDIR /data

# No ENTRYPOINT: init-containers and tests override `command` fully.
# Default CMD prints version table — useful for `docker run` discovery.
CMD ["/bin/sh", "-c", "for f in /data/*; do printf '=== %s ===\\n' \"$(basename \"$f\")\"; \"$f\" --version 2>&1 || true; done"]
