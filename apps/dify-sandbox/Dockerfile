# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

FROM langgenius/dify-sandbox:main@sha256:cb076f71cc84c14d4e4f7753ff95c4ba70a3b5816962b4f93bcf42f23a6e5cb8

ARG UID=5000
ARG GID=5000

# ---- mise (apt channel) --------------------------------------------------
RUN apt update -y && \
    install -dm 755 /etc/apt/keyrings && \
    curl -fSs https://mise.en.dev/gpg-key.pub | tee /etc/apt/keyrings/mise-archive-keyring.asc 1> /dev/null && \
    echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.asc] https://mise.en.dev/deb stable main" | tee /etc/apt/sources.list.d/mise.list && \
    apt update -y && apt install -y --no-install-recommends mise && \
    rm -rf /var/lib/apt/lists/*

# ---- non-root sandbox user for exec sessions -----------------------------
RUN groupadd -f -g "${GID}" sandbox && \
    useradd -m -u "${UID}" -g sandbox -s /bin/bash -d /home/sandbox sandbox 2>/dev/null; \
    usermod -s /bin/bash -d /home/sandbox sandbox 2>/dev/null; \
    mkdir -p /home/sandbox && chown -R sandbox:sandbox /home/sandbox

# ---- system-wide mise: tools available to root and sandbox ----------------
ENV MISE_DATA_DIR=/usr/local/mise \
    PATH=/usr/local/mise/shims:/usr/local/bin:/usr/bin:/bin

COPY app /app
RUN chmod +x /app/entrypoint.sh && \
    install -dm 755 /usr/local/mise && \
    install -dm 755 /home/sandbox/.config/mise && \
    install -dm 755 /root/.config/mise && \
    install -m 644 /app/.mise.toml /home/sandbox/.config/mise/config.toml && \
    ln -sf /home/sandbox/.config/mise/config.toml /root/.config/mise/config.toml

# install tools globally; mise shims on PATH for all users
RUN mise trust /home/sandbox/.config/mise/config.toml && \
    mise install && \
    mise reshim && \
    echo 'eval "$(mise activate bash)"' >> /home/sandbox/.bashrc && \
    echo 'eval "$(mise activate bash)"' >> /root/.bashrc

# override upstream entrypoint: activate mise, then run the sandbox server
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/main"]
