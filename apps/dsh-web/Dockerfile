# syntax=docker/dockerfile:1.7

# dsh-web — agent workbench appliance (nix-ops): base deepseek-harness + web
# profile (dsh-ui-auth + dshmarket), plugins baked at build, runtime keeps
# full install/update ability (pnpm stays; runtime `dsh plugin add` works and
# is covered by the container test contract).
# Build needs the base in context: bake target dsh-base builds
# ../deepseek-harness first (see docker-bake.hcl), or pass
# --build-context deepseek-harness=docker-image://ghcr.io/soulwhisper/deepseek-harness:latest

# Plugin pins are deliberately NOT renovate-tracked (coupled pins): the
# workbench has full runtime install/update ability, so plugin updates are
# user-driven inside the appliance — the image is not rebuilt for plugin
# sub-versions. Only upstreams are tracked: dsh core (via the base image).
ARG DSH_UI_AUTH_VERSION=0.7.1
ARG DSH_MARKET_VERSION=1.66.44

FROM deepseek-harness

ARG DSH_UI_AUTH_VERSION
ARG DSH_MARKET_VERSION

USER root
RUN mkdir -p /opt/dsh-seed && chown node:node /opt/dsh-seed

USER node
# Bake the web profile into the seed template. dsh plugin add writes into
# $DSH_HOME/profiles/web; keep it out of the runtime home.
RUN export DSH_HOME=/opt/dsh-seed HOME=/tmp; \
    set -eux; \
    dsh plugin --profile web add "dsh-ui-auth@${DSH_UI_AUTH_VERSION}"; \
    dsh plugin --profile web add "dshmarket@${DSH_MARKET_VERSION}"; \
    test -d "${DSH_HOME}/profiles/web/node_modules/dsh-ui-auth"; \
    test -d "${DSH_HOME}/profiles/web/node_modules/dshmarket"; \
    { \
      echo "dsh-ui-auth=${DSH_UI_AUTH_VERSION}"; \
      echo "dshmarket=${DSH_MARKET_VERSION}"; \
    } > "${DSH_HOME}/BUNDLED_VERSIONS"; \
    # dsh plugin add writes some files root-only; the entrypoint seeds as node
    chmod -R a+rX "${DSH_HOME}"

ENV DSH_PROFILE=web

ARG IMAGE_VERSION=0.2.1-rc.2-r1
ARG IMAGE_REVISION=unknown
LABEL org.opencontainers.image.title="dsh-web" \
      org.opencontainers.image.description="DeepSeek Harness agent workbench: web profile with dsh-ui-auth + dshmarket baked, runtime plugin installs enabled" \
      org.opencontainers.image.vendor="soulwhisper" \
      org.opencontainers.image.source="https://github.com/soulwhisper/containers" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}" \
      io.github.soulwhisper.dsh-web.plugins="dsh-ui-auth@${DSH_UI_AUTH_VERSION},dshmarket@${DSH_MARKET_VERSION}"
