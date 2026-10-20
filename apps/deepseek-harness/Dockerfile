# syntax=docker/dockerfile:1.7

# deepseek-harness — DeepSeek Harness (@deepseek-ai/dsh) appliance image with
# two switchable profiles:
#   web    — agent workbench (nix-ops): dsh-ui-auth + dsh-pocket-relay + dshmarket
#   tavern — RP chat (cluster):       flizzywine/dsh-tavern + dsh-ui-auth + dsh-pocket
# Profile is selected at runtime by DSH_PROFILE (default: web) or $1.
# All plugins are baked at build time into /opt/dsh-seed; the entrypoint
# seeds them into $DSH_HOME (PVC) without touching user data
# (profile-data/). No upstream fork: dsh comes from npm, plugins from
# npm/github, version-pinned with build-time assertions.

ARG NODE_IMAGE=docker.io/library/node:24-trixie
# renovate: datasource=npm depName=@deepseek-ai/dsh
ARG DSH_VERSION=0.2.1-rc.2
# renovate: datasource=npm depName=pnpm
ARG PNPM_VERSION=10.15.1
# renovate: datasource=npm depName=dsh-ui-auth
ARG DSH_UI_AUTH_VERSION=0.7.1
# renovate: datasource=npm depName=dsh-pocket
ARG DSH_POCKET_VERSION=2.12.0
# renovate: datasource=npm depName=dsh-pocket-relay
ARG DSH_POCKET_RELAY_VERSION=1.0.2
# renovate: datasource=npm depName=dshmarket
ARG DSH_MARKET_VERSION=1.66.44
# renovate: datasource=github-tags depName=flizzywine/dsh-tavern
ARG DSH_TAVERN_REF=v3.0

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

# ---- Plugin baking ---------------------------------------------------------
# `dsh plugin add` writes into $DSH_HOME/profiles/<profile>; bake both
# profiles into a throwaway DSH_HOME that ships as the seed template.
FROM installer AS plugins

ARG DSH_UI_AUTH_VERSION
ARG DSH_POCKET_VERSION
ARG DSH_POCKET_RELAY_VERSION
ARG DSH_MARKET_VERSION
ARG DSH_TAVERN_REF

ENV DSH_HOME=/opt/dsh-seed \
    HOME=/tmp

RUN set -eux; \
    # web profile: auth gateway + self-hosted phone relay + plugin market
    dsh plugin --profile web add "dsh-ui-auth@${DSH_UI_AUTH_VERSION}"; \
    dsh plugin --profile web add "dsh-pocket-relay@${DSH_POCKET_RELAY_VERSION}"; \
    dsh plugin --profile web add "dshmarket@${DSH_MARKET_VERSION}"; \
    # tavern profile: tavern runtime + auth gateway + phone access
    dsh plugin --profile tavern add "github:flizzywine/dsh-tavern#${DSH_TAVERN_REF}"; \
    dsh plugin --profile tavern add "dsh-ui-auth@${DSH_UI_AUTH_VERSION}"; \
    dsh plugin --profile tavern add "dsh-pocket@${DSH_POCKET_VERSION}"; \
    # build-time assertions: every baked plugin is on disk
    test -d "${DSH_HOME}/profiles/web/node_modules/dsh-ui-auth"; \
    test -d "${DSH_HOME}/profiles/web/node_modules/dsh-pocket-relay"; \
    test -d "${DSH_HOME}/profiles/web/node_modules/dshmarket"; \
    test -d "${DSH_HOME}/profiles/tavern/node_modules/dsh-ui-auth"; \
    test -d "${DSH_HOME}/profiles/tavern/node_modules/dsh-pocket"; \
    { \
      echo "dsh=${DSH_VERSION}"; \
      echo "pnpm=${PNPM_VERSION}"; \
      echo "dsh-ui-auth=${DSH_UI_AUTH_VERSION}"; \
      echo "dsh-pocket=${DSH_POCKET_VERSION}"; \
      echo "dsh-pocket-relay=${DSH_POCKET_RELAY_VERSION}"; \
      echo "dshmarket=${DSH_MARKET_VERSION}"; \
      echo "dsh-tavern=${DSH_TAVERN_REF}"; \
    } > "${DSH_HOME}/BUNDLED_VERSIONS"; \
    # dsh plugin add writes some files root-only (package.json, .plugin-manager);
    # the runtime entrypoint seeds $DSH_HOME as the node user, so the seed
    # template must be world-readable (no credentials exist at build time)
    chmod -R a+rX "${DSH_HOME}"

# ---- Runtime ---------------------------------------------------------------
FROM ${NODE_IMAGE}

ARG NODE_IMAGE
ARG DSH_VERSION
ARG PNPM_VERSION
ARG DSH_UI_AUTH_VERSION
ARG DSH_POCKET_VERSION
ARG DSH_POCKET_RELAY_VERSION
ARG DSH_MARKET_VERSION
ARG DSH_TAVERN_REF

# Agent-facing CLI tools only — no browser/X11 stack (that's the heavyweight
# variant; add a `desktop` build target if browser-use is ever needed).
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

COPY --from=plugins /usr/local/lib/node_modules/@deepseek-ai/dsh /usr/local/lib/node_modules/@deepseek-ai/dsh
COPY --from=plugins /usr/local/lib/node_modules/pnpm /usr/local/lib/node_modules/pnpm
COPY --from=plugins /opt/dsh-seed /opt/dsh-seed
COPY app/entrypoint.sh /app/entrypoint.sh
COPY app/dsh-web-relay.cjs /app/dsh-web-relay.cjs

RUN chmod 0755 /app/entrypoint.sh \
    && ln -s ../lib/node_modules/@deepseek-ai/dsh/lib/bin.js /usr/local/bin/dsh \
    && ln -s ../lib/node_modules/pnpm/bin/pnpm.cjs /usr/local/bin/pnpm \
    && ln -s ../lib/node_modules/pnpm/bin/pnpx.cjs /usr/local/bin/pnpx \
    && test "$(dsh --version)" = "${DSH_VERSION}" \
    && test "$(pnpm --version)" = "${PNPM_VERSION}" \
    && test -f /opt/dsh-seed/BUNDLED_VERSIONS

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

ARG IMAGE_VERSION=0.2.1-rc.2-r0
ARG IMAGE_REVISION=unknown
LABEL org.opencontainers.image.title="deepseek-harness" \
      org.opencontainers.image.description="DeepSeek Harness appliance: web (agent workbench) and tavern (RP chat) profiles, plugins baked at build time" \
      org.opencontainers.image.vendor="soulwhisper" \
      org.opencontainers.image.source="https://github.com/soulwhisper/containers" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}" \
      io.github.soulwhisper.deepseek-harness.dsh.version="${DSH_VERSION}" \
      io.github.soulwhisper.deepseek-harness.tavern.ref="${DSH_TAVERN_REF}" \
      io.github.soulwhisper.deepseek-harness.plugins="dsh-ui-auth@${DSH_UI_AUTH_VERSION},dsh-pocket@${DSH_POCKET_VERSION},dsh-pocket-relay@${DSH_POCKET_RELAY_VERSION},dshmarket@${DSH_MARKET_VERSION}" \
      io.github.soulwhisper.deepseek-harness.upstream.npm="@deepseek-ai/dsh@${DSH_VERSION}"

ENTRYPOINT ["/usr/bin/tini", "--", "/app/entrypoint.sh"]
CMD []
