# Build stage — fetch source, install deps, compile Next.js standalone output.
# VERSION is a git commit SHA (no tags/releases; renovate tracks git-refs).
ARG VERSION=33085b6608fe7174e0fb75e46c220348863b1c19

FROM node:22-alpine AS builder
ARG VERSION
WORKDIR /app

# Fetch source tarball at the pinned commit (no git clone — lighter).
ADD "https://github.com/langgenius/webapp-conversation/archive/${VERSION}.tar.gz" /tmp/repo.tar.gz
RUN tar xzf /tmp/repo.tar.gz -C /tmp && \
    cp -r /tmp/webapp-conversation-*/. /app/ && \
    rm -rf /tmp/repo.tar.gz /tmp/webapp-conversation-*

# ── Locale patches: Chinese (zh-Hans) as default ──
RUN sed -i "s/defaultLocale: 'en'/defaultLocale: 'zh-Hans'/" i18n/index.ts
RUN sed -i "s/default_language: 'en'/default_language: 'zh-Hans'/" config/index.ts
RUN sed -i "s/lng: 'en'/lng: 'zh-Hans'/" i18n/i18next-config.ts
RUN sed -i "s/fallbackLng: 'en'/fallbackLng: 'zh-Hans'/" i18n/i18next-config.ts
RUN sed -i "s/fallbackLng: 'en'/fallbackLng: 'zh-Hans'/" i18n/i18next-serverside-config.ts

# Install deps and build. output: 'standalone' produces .next/standalone/server.js.
RUN yarn install --frozen-lockfile
RUN yarn build

# ── Runner stage ──
FROM node:22-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001 -G nodejs

COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

CMD ["node", "server.js"]
