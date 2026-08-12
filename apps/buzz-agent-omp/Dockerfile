# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# ---- buzz agent body (sprig multicall: buzz-acp/buzz-agent/buzz-dev-mcp/buzz
# CLI/rg/tree/git-{credential,sign}-nostr), statically linked upstream --------
FROM ghcr.io/block/buzz-sprig:main@sha256:e74d01228d98579d410465ec46c58fa94d20dae0914fc0ffb25ce102fc89531d AS sprig

FROM debian:13-slim@sha256:3a39a0592364683e6bab97937b72cad5a8fa6dcbbee90edb3bb48c7f8e94f258

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
RUN groupadd -g "${GID}" agent && \
    useradd -m -u "${UID}" -g "${GID}" -s /bin/bash agent && \
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
ENV HOME=/home/agent \
    MISE_DATA_DIR=/home/agent/.local/share/mise \
    PATH=/home/agent/.local/share/mise/shims:/usr/local/bin:/usr/bin:/bin \
    BUZZ_ACP_AGENT_COMMAND=omp \
    BUZZ_ACP_AGENT_ARGS=acp

COPY app /app
RUN chmod +x /app/entrypoint.sh && \
    install -dm 755 -o agent -g agent /home/agent/.config/mise && \
    install -m 644 -o agent -g agent /app/.mise.toml /home/agent/.config/mise/config.toml

USER 2000:2000
WORKDIR /home/agent

# install every tool from the global .mise.toml; expose them via shims on PATH
RUN mise trust /home/agent/.config/mise/config.toml && \
    mise install && \
    mise reshim && \
    echo 'eval "$(mise activate bash)"' >> /home/agent/.bashrc

ENTRYPOINT ["/app/entrypoint.sh"]
