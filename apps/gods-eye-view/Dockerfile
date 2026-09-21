# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# God's Eye View — live open-source spatial intelligence on a 3D globe.
#   https://github.com/bilawalsidhu/gods-eye-view
#
# Upstream ships no container image and no standalone production server: the
# supported runtime is a Vite dev server that serves the client and brokers
# every third-party API same-origin under /api/* (server/providers/*), exactly
# like upstream's own headless launcher (scripts/pinokio-start.mjs). This
# image mirrors that launcher: full checkout + node_modules, keys via env.
#
# HOST=0.0.0.0 is what upstream's vite config keys on to set allowedHosts:
# true — required so a reverse proxy (kgateway/authentik forward-auth) can
# forward the public Host header.
#
# Client-exposed keys (GOOGLE_MAPS_API_KEY, CESIUM_ION_TOKEN) are injected
# into the browser bundle by design (vite `define` at server start). All
# other keys stay server-side inside the /api/* proxies.
#
# Service port:
#   4173 — web UI + /api/* provider proxies

ARG VERSION=main

FROM node:24-bookworm-slim@sha256:0e0ff40c39bc087845bfb27465a0df4ea419520094bc35842ff83dd8cbe6f9b6 AS upstream
ARG VERSION
RUN apt-get -qq update \
    && apt-get -qq install --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN git init -q . \
    && git remote add origin https://github.com/bilawalsidhu/gods-eye-view.git \
    && git fetch -q --depth 1 origin "${VERSION}" \
    && git checkout -q FETCH_HEAD \
    && rm -rf .git

# Dependencies in their own layer: puppeteer (QA-only dep) downloads Chromium
# on postinstall — skip it; sharp/rollup/esbuild pull platform-correct prebuilt
# binaries as npm optional deps.
FROM node:24-bookworm-slim@sha256:0e0ff40c39bc087845bfb27465a0df4ea419520094bc35842ff83dd8cbe6f9b6 AS deps
WORKDIR /app
COPY --from=upstream /app/package.json /app/package-lock.json ./
ENV PUPPETEER_SKIP_DOWNLOAD=true
RUN npm ci --no-audit --no-fund

FROM node:24-bookworm-slim@sha256:0e0ff40c39bc087845bfb27465a0df4ea419520094bc35842ff83dd8cbe6f9b6
ENV HOST=0.0.0.0 \
    PORT=4173 \
    # Upstream recommends enabling these whenever HOST is not localhost: the
    # server brokers spend-quota keys. In-cluster auth is forward-auth, but
    # the limits also cap accidental spend from any direct port-forward.
    GEV_RATELIMIT_OPENAI_PER_MIN=30 \
    GEV_RATELIMIT_GOOGLE_PER_MIN=60
WORKDIR /app
COPY --from=upstream /app .
COPY --from=deps /app/node_modules ./node_modules
RUN chown -R node:node /app
USER node
EXPOSE 4173
CMD ["node", "node_modules/vite/bin/vite.js"]
