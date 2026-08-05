# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

# ---- buzz agent body (sprig multicall: buzz-acp/buzz-agent/buzz-dev-mcp/buzz
# CLI/rg/tree/git-{credential,sign}-nostr), statically linked upstream --------
FROM ghcr.io/block/buzz-sprig:main@sha256:135a80c2dd87ebe9923de707dc7f086100d495e3878909e1426947a5f0244606 AS sprig

FROM debian:13-slim@sha256:020c0d20b9880058cbe785a9db107156c3c75c2ac944a6aa7ab59f2add76a7bd

ARG UID=2000
ARG GID=2000

# ---- base OS deps + mise (apt channel, same as sibling images) ------------
RUN apt update -y && \
    apt install -y --no-install-recommends \
        bash ca-certificates curl git openssh-client && \
    install -dm 755 /etc/apt/keyrings && \
    curl -fSs https://mise.en.dev/gpg-key.pub | tee /etc/apt/keyrings/mise-archive-keyring.asc 1> /dev/null && \
    echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.asc] https://mise.en.dev/deb stable main" | tee /etc/apt/sources.list.d/mise.list && \
    apt update -y && apt install -y --no-install-recommends mise && \
    rm -rf /var/lib/apt/lists/*

# ---- non-root agent user (uid=gid=2000, bash login shell) ------------------
# git: wildcard safe.directory for mounted workspaces; x509 commit signing via
# the agent's nostr key, mirroring the upstream sprig image contract.
RUN groupadd -g "${GID}" dev && \
    useradd -m -u "${UID}" -g "${GID}" -s /bin/bash dev && \
    git config --system --add safe.directory '*' && \
    git config --system gpg.format x509 && \
    git config --system gpg.x509.program /usr/local/bin/git-sign-nostr && \
    git config --system commit.gpgSign true && \
    git config --system tag.gpgSign true

# ---- buzz sprig multicall binaries -----------------------------------------
COPY --from=sprig /usr/local/bin/sprig /usr/local/bin/sprig
RUN for name in \
        buzz-acp buzz-agent buzz-dev-mcp rg tree buzz \
        git-credential-nostr git-sign-nostr; do \
        ln -s sprig "/usr/local/bin/${name}"; \
    done

# ---- env: mise shims on PATH; omp as the default ACP agent -----------------
# Identity, relay and model routing are supplied entirely by container env at
# deploy time (from Secrets): BUZZ_PRIVATE_KEY, BUZZ_RELAY_URL, BUZZ_AUTH_TAG,
# BUZZ_API_TOKEN, and omp provider credentials. Only harness defaults are baked.
ENV HOME=/home/dev \
    MISE_DATA_DIR=/home/dev/.local/share/mise \
    PATH=/home/dev/.local/share/mise/shims:/usr/local/bin:/usr/bin:/bin \
    BUZZ_ACP_AGENT_COMMAND=omp \
    BUZZ_ACP_AGENT_ARGS=acp

COPY app /app
RUN chmod +x /app/entrypoint.sh && \
    install -dm 755 -o dev -g dev /home/dev/.config/mise && \
    install -m 644 -o dev -g dev /app/.mise.toml /home/dev/.config/mise/config.toml

USER 2000:2000
WORKDIR /home/dev

# install every tool from the global .mise.toml; expose them via shims on PATH
RUN mise trust /home/dev/.config/mise/config.toml && \
    mise install && \
    mise reshim && \
    echo 'eval "$(mise activate bash)"' >> /home/dev/.bashrc

ENTRYPOINT ["/app/entrypoint.sh"]
