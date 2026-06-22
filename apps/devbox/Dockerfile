# syntax=docker/dockerfile:1

FROM debian:13

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

# ---- non-root devbox user (uid=gid=2000, bash login shell) ----------------
RUN groupadd -g "${GID}" dev && \
    useradd -m -u "${UID}" -g "${GID}" -s /bin/bash dev && \
    git config --system --add safe.directory '*'

# ---- env: mise shims on PATH; provider-neutral by design -------------------
# NO endpoint or model is baked, so the image never forces DeepSeek (or any
# backend). Routing is supplied entirely by container env at deploy time:
#   ANTHROPIC_BASE_URL    -> Anthropic-compatible route
#   ANTHROPIC_AUTH_TOKEN  -> bearer token (from a Secret)
#   ANTHROPIC_MODEL / ANTHROPIC_DEFAULT_{OPUS,SONNET,HAIKU}_MODEL /
#   CLAUDE_CODE_SUBAGENT_MODEL  -> model ids
# Only provider-neutral hygiene flags are kept here.
ENV HOME=/home/dev \
    MISE_DATA_DIR=/home/dev/.local/share/mise \
    PATH=/home/dev/.local/share/mise/shims:/usr/local/bin:/usr/bin:/bin \
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

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

# No fixed command: idle so you can `kubectl exec -it -- bash`
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["sleep", "infinity"]
