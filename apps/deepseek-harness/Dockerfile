# syntax=docker/dockerfile:1.7

# deepseek-harness — BASEMENT image: node + DeepSeek Harness (@deepseek-ai/dsh)
# + container fixes (web relay, seeding entrypoint). NO plugins, NO profiles —
# those are child-image concerns (dsh-web / dsh-tavern), because dsh plugins
# activate per-profile, not globally. The entrypoint and seed contract are
# shared: children bake /opt/dsh-seed and this same entrypoint seeds it into
# $DSH_HOME on first boot.

ARG NODE_IMAGE=docker.io/library/node:24-trixie
# renovate: datasource=npm depName=@deepseek-ai/dsh
ARG DSH_VERSION=0.2.1-rc.2
# renovate: datasource=npm depName=pnpm
ARG PNPM_VERSION=11.25.0

FROM ${NODE_IMAGE} AS installer

ARG DSH_VERSION
ARG PNPM_VERSION

# node-pty publishes prebuilds only for some Linux archs; keep the toolchain
# for node-gyp fallback. allow-scripts is a whitelist for packages whose
# postinstall is required (native modules, genai protobuf codegen).
RUN sed -i 's|http://deb.debian.org|https://deb.debian.org|g' /etc/apt/sources.list.d/debian.sources \
    && apt-get update \
    && apt-get install --yes --no-install-recommends \
      build-essential \
      ca-certificates \
      python3 \
    && rm -rf /var/lib/apt/lists/* \
    && npm install --global --omit=dev --no-audit --no-fund \
      --allow-scripts=@deepseek-ai/dsh-subprocess-local,koffi,node-pty,@google/genai,protobufjs \
      "@deepseek-ai/dsh@${DSH_VERSION}" \
      "pnpm@${PNPM_VERSION}" \
    && test "$(dsh --version)" = "${DSH_VERSION}" \
    && test "$(pnpm --version)" = "${PNPM_VERSION}" \
    && npm cache clean --force

# ---- Runtime ---------------------------------------------------------------
FROM ${NODE_IMAGE}

ARG NODE_IMAGE
ARG DSH_VERSION
ARG PNPM_VERSION

# Agent-facing CLI tools. No browser/X11 stack (that's a future `desktop`
# build target). Full (non-slim) base is deliberate: children need the
# toolchain for native plugin builds and runtime plugin installs.
RUN sed -i 's|http://deb.debian.org|https://deb.debian.org|g' /etc/apt/sources.list.d/debian.sources \
    && apt-get update \
    && apt-get install --yes --no-install-recommends \
      ca-certificates \
      curl \
      git \
      jq \
      less \
      openssh-client \
      procps \
      ripgrep \
      rsync \
      tini \
      unzip \
      zip \
    && rm -rf /var/lib/apt/lists/*

COPY --from=installer /usr/local/lib/node_modules/@deepseek-ai/dsh /usr/local/lib/node_modules/@deepseek-ai/dsh
COPY --from=installer /usr/local/lib/node_modules/pnpm /usr/local/lib/node_modules/pnpm
COPY app/entrypoint.sh /app/entrypoint.sh
COPY app/dsh-web-relay.cjs /app/dsh-web-relay.cjs

RUN chmod 0755 /app/entrypoint.sh \
    && chmod 0755 /usr/local/lib/node_modules/pnpm/bin/pnpm.cjs \
      /usr/local/lib/node_modules/pnpm/bin/pnpx.cjs \
    && ln -s ../lib/node_modules/@deepseek-ai/dsh/lib/bin.js /usr/local/bin/dsh \
    && ln -s ../lib/node_modules/pnpm/bin/pnpm.cjs /usr/local/bin/pnpm \
    && ln -s ../lib/node_modules/pnpm/bin/pnpx.cjs /usr/local/bin/pnpx \
    && test "$(dsh --version)" = "${DSH_VERSION}" \
    && test "$(pnpm --version)" = "${PNPM_VERSION}"

ENV DSH_HOME=/home/node/.dsh \
    DSH_TELEMETRY_DISABLED=1 \
    DSH_PROFILE=web \
    DSH_PORT=3080 \
    HOME=/home/node \
    NPM_CONFIG_CACHE=/home/node/.dsh/npm-cache \
    XDG_CACHE_HOME=/tmp/.cache \
    XDG_CONFIG_HOME=/tmp/.config \
    XDG_DATA_HOME=/tmp/.local/share

RUN mkdir -p "${DSH_HOME}" && chown -R node:node "${DSH_HOME}"

USER node
WORKDIR /workspace

EXPOSE 3080

ARG IMAGE_VERSION=0.2.1-rc.2-r1
ARG IMAGE_REVISION=unknown
LABEL org.opencontainers.image.title="deepseek-harness" \
      org.opencontainers.image.description="DeepSeek Harness basement image: node + dsh + container fixes (no plugins). Parent of dsh-web / dsh-tavern appliances." \
      org.opencontainers.image.vendor="soulwhisper" \
      org.opencontainers.image.source="https://github.com/soulwhisper/containers" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}" \
      io.github.soulwhisper.deepseek-harness.dsh.version="${DSH_VERSION}" \
      io.github.soulwhisper.deepseek-harness.pnpm.version="${PNPM_VERSION}" \
      io.github.soulwhisper.deepseek-harness.upstream.npm="@deepseek-ai/dsh@${DSH_VERSION}"

ENTRYPOINT ["/usr/bin/tini", "--", "/app/entrypoint.sh"]
CMD []
