# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

# Build stage — fetch source, install deps, compile Next.js standalone output.
# VERSION is a git commit SHA (no tags/releases; renovate tracks git-refs).
ARG VERSION=33085b6608fe7174e0fb75e46c220348863b1c19

FROM node:22-alpine@sha256:16e22a550f3863206a3f701448c45f7912c6896a62de43add43bb9c86130c3e2 AS builder
ARG VERSION
WORKDIR /app

# Fetch source tarball at the pinned commit (no git clone — lighter).
ADD "https://github.com/langgenius/webapp-conversation/archive/${VERSION}.tar.gz" /tmp/repo.tar.gz
RUN tar xzf /tmp/repo.tar.gz -C /tmp && \
    cp -r /tmp/webapp-conversation-*/. /app/ && \
    rm -rf /tmp/repo.tar.gz /tmp/webapp-conversation-*

# Install deps and build. output: 'standalone' produces .next/standalone/server.js.
RUN yarn install --frozen-lockfile
RUN yarn build

# Runner stage — minimal production image with non-root user.
FROM node:22-alpine@sha256:16e22a550f3863206a3f701448c45f7912c6896a62de43add43bb9c86130c3e2 AS runner
WORKDIR /app

ENV NODE_ENV=production

# Non-root user matching Next.js conventions (1001:1001).
RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001 -G nodejs

COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

CMD ["node", "server.js"]
