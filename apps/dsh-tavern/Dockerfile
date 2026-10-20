# syntax=docker/dockerfile:1.7

# dsh-tavern — RP chat appliance (cluster): base deepseek-harness + the
# flizzywine/dsh-tavern profile installed via the VERIFIED CUSTOM RECIPE
# (not install.sh — that path brings a second, private dsh runtime).
# Recipe (every step build-asserted):
#   1. clone repo at pinned ref            → /opt/dsh-tavern/apps/dsh-tavern
#   2. pnpm install --frozen-lockfile      → tavern-plugin workspace deps
#      (this is what plain `dsh plugin add` misses — dsh-tavern#300)
#   3. materialize profile: generated package.json (link: deps), PRUNED
#      pnpm-workspace.yaml (modeled on the installer's output, not the
#      source's), patches/ trimmed to the used set, cordis.patch.yml
#   4. pnpm install in the profile
#   5. dsh plugin add: eco plugins
#   6. host UX patches (font limit, chat answer step) applied to THIS image's
#      dsh tree only
#   7. freeze: pnpm binaries removed — runtime plugin changes are impossible
#      by design; updates happen only via image rebuilds.
# Build needs the base in context (see docker-bake.hcl).

# renovate: datasource=github-tags depName=flizzywine/dsh-tavern
ARG DSH_TAVERN_REF=v3.0
# Eco-plugin pins are deliberately NOT renovate-tracked (coupled pins): they
# are re-evaluated only when DSH_TAVERN_REF drifts upstream — rebuilding the
# appliance for an eco-plugin's own release cycle is unwanted.
ARG DSH_WRONGBOOK_REF=v1.33.0
ARG DSH_TAVERN_IMAGE_REF=v2.9.0
ARG DSH_KNOCK_KNOCK_REF=v0.7.1
# no release tags on the current line — pinned commit; bump manually
ARG DSH_NARRATIVE_ANCHOR_SHA=2ee3da047e5fd6056aef361744655be99b4f7713

FROM deepseek-harness

ARG DSH_TAVERN_REF
ARG DSH_WRONGBOOK_REF
ARG DSH_TAVERN_IMAGE_REF
ARG DSH_KNOCK_KNOCK_REF
ARG DSH_NARRATIVE_ANCHOR_SHA
# filled in from the base image build; used in the profile manifest
ARG DSH_VERSION=0.2.1-rc.2

ENV DSH_TAVERN_SRC=/opt/dsh-tavern/apps/dsh-tavern \
    DSH_PROFILE=tavern

USER root

# 1-2. source + workspace deps (as node: keeps files readable for the seed)
RUN export HOME=/tmp; \
    set -eux; \
    mkdir -p /opt/dsh-tavern/apps /opt/dsh-seed/profiles/tavern; \
    git clone --depth 1 --branch "${DSH_TAVERN_REF}" \
      https://github.com/flizzywine/dsh-tavern "${DSH_TAVERN_SRC}"; \
    cd "${DSH_TAVERN_SRC}"; \
    pnpm install --frozen-lockfile; \
    test -d node_modules/jsdom; \
    chown -R node:node /opt/dsh-tavern /opt/dsh-seed

# 3-5. profile materialization + plugin adds (as node, like the runtime)
USER node
COPY --chown=node:node app/profile-package.json.tmpl /tmp/profile-package.json.tmpl
COPY --chown=node:node app/profile-pnpm-workspace.yaml /tmp/profile-pnpm-workspace.yaml
RUN export DSH_HOME=/opt/dsh-seed HOME=/tmp; \
    set -eux; \
    PROF="${DSH_HOME}/profiles/tavern"; \
    sed -e "s|@SRC@|${DSH_TAVERN_SRC}|g" -e "s|@DSH_VERSION@|${DSH_VERSION}|g" \
      /tmp/profile-package.json.tmpl > "${PROF}/package.json"; \
    cp /tmp/profile-pnpm-workspace.yaml "${PROF}/pnpm-workspace.yaml"; \
    cp "${DSH_TAVERN_SRC}/cordis.patch.yml" "${PROF}/cordis.patch.yml"; \
    mkdir -p "${PROF}/patches"; \
    cp "${DSH_TAVERN_SRC}/patches/dsh-better-sidebar@0.28.0.patch" "${PROF}/patches/"; \
    cd "${PROF}"; \
    pnpm install; \
    test -d "${PROF}/node_modules/dsh-tavern-plugin"; \
    dsh plugin --profile tavern add "github:XGUIMAX/dsh-tavern-wrongbook#${DSH_WRONGBOOK_REF}"; \
    dsh plugin --profile tavern add "github:weixinlll/dsh-tavern-image#${DSH_TAVERN_IMAGE_REF}"; \
    dsh plugin --profile tavern add "github:dearacknowl/dsh-tavern-knock-knock#${DSH_KNOCK_KNOCK_REF}"; \
    dsh plugin --profile tavern add "github:ActivityTime/narrative-anchor#${DSH_NARRATIVE_ANCHOR_SHA}"; \
    for p in dsh-tavern-wrongbook dsh-tavern-comfy dsh-tavern-knock-knock narrative-anchor; do \
      test -d "${PROF}/node_modules/$p"; \
    done

# 6. host UX patches (idempotent; skip-safe if a future dsh changes shape).
# chat-answer-step exists only on main (post-v3.0) — applied automatically
# once a tagged ref ships it; font-limit exists since v3.0.
USER root
RUN node --input-type=module -e "\
import { existsSync } from 'node:fs'; \
const src = process.env.DSH_TAVERN_SRC + '/bin'; \
const tree = '/usr/local/lib/node_modules/@deepseek-ai/dsh/node_modules/@deepseek-ai'; \
const { patchThemeFontLimit } = await import(src + '/theme-font-limit.mjs'); \
console.log('font-limit patch:', patchThemeFontLimit(tree + '/dsh-client-ui-theme')); \
if (existsSync(src + '/chat-answer-step-patch.mjs')) { \
  const { patchChatAnswerStep } = await import(src + '/chat-answer-step-patch.mjs'); \
  console.log('chat-answer-step patch:', patchChatAnswerStep(tree + '/dsh-client-ui-chat')); \
} else { \
  console.log('chat-answer-step patch: script not in this ref, skipped'); \
}"

# 7. version stamp, permissions, freeze
RUN export DSH_HOME=/opt/dsh-seed; \
    set -eux; \
    { \
      echo "dsh-tavern=${DSH_TAVERN_REF}"; \
      echo "dsh-tavern-wrongbook=${DSH_WRONGBOOK_REF}"; \
      echo "dsh-tavern-image=${DSH_TAVERN_IMAGE_REF}"; \
      echo "dsh-tavern-knock-knock=${DSH_KNOCK_KNOCK_REF}"; \
      echo "narrative-anchor=${DSH_NARRATIVE_ANCHOR_SHA}"; \
    } > "${DSH_HOME}/BUNDLED_VERSIONS"; \
    chmod -R a+rX /opt/dsh-seed /opt/dsh-tavern; \
    chown -R node:node /opt/dsh-seed; \
    # pnpm writes link: deps as RELATIVE symlinks anchored at the seed root
    # (/opt/dsh-seed) — they dangle wherever the entrypoint copies the
    # profile ($DSH_HOME, tests, upgrades). Rewrite the three link: deps as
    # ABSOLUTE symlinks so they survive any copy location.
    for d in dsh-tavern-plugin dsh-dream-skin dsh-tavern-remote; do \
      ln -sfn "$(realpath -m "${DSH_HOME}/profiles/tavern/node_modules/$d")" \
        "${DSH_HOME}/profiles/tavern/node_modules/$d"; \
    done; \
    # freeze: no runtime plugin installs/updates (stability contract)
    rm -f /usr/local/bin/pnpm /usr/local/bin/pnpx

USER node

ARG IMAGE_VERSION=0.2.1-rc.2-r1
ARG IMAGE_REVISION=unknown
LABEL org.opencontainers.image.title="dsh-tavern" \
      org.opencontainers.image.description="DeepSeek Harness RP chat appliance: dsh-tavern profile + ecosystem plugins, frozen (no runtime plugin changes); auth via the authentik portal" \
      org.opencontainers.image.vendor="soulwhisper" \
      org.opencontainers.image.source="https://github.com/soulwhisper/containers" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}" \
      io.github.soulwhisper.dsh-tavern.ref="${DSH_TAVERN_REF}" \
      io.github.soulwhisper.dsh-tavern.plugins="wrongbook@${DSH_WRONGBOOK_REF},tavern-image@${DSH_TAVERN_IMAGE_REF},knock-knock@${DSH_KNOCK_KNOCK_REF},narrative-anchor@${DSH_NARRATIVE_ANCHOR_SHA}"
