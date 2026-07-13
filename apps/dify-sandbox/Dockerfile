# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

FROM debian:13-slim@sha256:28de0877c2189802884ccd20f15ee41c203573bd87bb6b883f5f46362d24c5c2

ARG UID=2000
ARG GID=2000

# ---- base OS deps + mise (apt channel) ----------------------------------
RUN apt update -y && \
    apt install -y --no-install-recommends \
        bash ca-certificates curl git libatomic1 && \
    install -dm 755 /etc/apt/keyrings && \
    curl -fSs https://mise.en.dev/gpg-key.pub | tee /etc/apt/keyrings/mise-archive-keyring.asc 1> /dev/null && \
    echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.asc] https://mise.en.dev/deb stable main" | tee /etc/apt/sources.list.d/mise.list && \
    apt update -y && apt install -y --no-install-recommends mise && \
    rm -rf /var/lib/apt/lists/*

# ---- non-root sandbox user (uid=gid=2000, bash login shell) -------------
RUN groupadd -g "${GID}" sandbox && \
    useradd -m -u "${UID}" -g "${GID}" -s /bin/bash sandbox

# ---- env: mise shims on PATH --------------------------------------------
ENV HOME=/home/sandbox \
    MISE_DATA_DIR=/home/sandbox/.local/share/mise \
    PATH=/home/sandbox/.local/share/mise/shims:/usr/local/bin:/usr/bin:/bin

COPY app /app
RUN chmod +x /app/entrypoint.sh && \
    install -dm 755 -o sandbox -g sandbox /home/sandbox/.config/mise && \
    install -m 644 -o sandbox -g sandbox /app/.mise.toml /home/sandbox/.config/mise/config.toml

USER 2000:2000
WORKDIR /home/sandbox

# install every tool from the global .mise.toml; expose them via shims on PATH
RUN mise trust /home/sandbox/.config/mise/config.toml && \
    mise install && \
    mise reshim && \
    echo 'eval "$(mise activate bash)"' >> /home/sandbox/.bashrc

# No fixed command: idle so you can `kubectl exec -it -- bash`
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["sleep", "infinity"]
