# deepseek-harness (base)

Basement image for the DeepSeek Harness appliances: **node + @deepseek-ai/dsh
+ container fixes, zero plugins**. Children (`dsh-web`, `dsh-tavern`) bake
their profiles on top — dsh plugins activate per-profile, so bundling plugins
here would activate nothing.

| Pin | Value | Notes |
| --- | ----- | ----- |
| `node` | `24-trixie` (full, deliberate) | toolchain for native plugin builds and runtime plugin installs |
| `@deepseek-ai/dsh` | `0.2.1-rc.2` | the only compatibility-critical track; renovate (npm) |
| `pnpm` | `11.25.0` | reproducibility pin; dsh has no pnpm constraint (no packageManager/engines). Chosen over runzhliu's 10.15.1 because the dsh-tavern repo lockfile is pnpm-11-format and every plugin operation used here was verified under 11 |

Fixes carried (design adopted from runzhliu/deepseek-harness-docker, MIT,
rewritten in-repo — no third-party image dependency):

- build-time version assertions (`dsh --version`, `pnpm --version`)
- npm hygiene: `--omit=dev --no-audit --no-fund` + postinstall allowlist
- `app/dsh-web-relay.cjs`: binds the container's concrete IPv4 → dsh's
  loopback-only web; multi-NIC ambiguity requires `DSH_WEB_RELAY_INTERFACE`
- seeding entrypoint (shared with children): if `/opt/dsh-seed` exists, a
  fresh `$DSH_HOME` receives it whole; on image bumps only plugin **code**
  (`profiles/*/node_modules`) is re-synced — `profile-data/` and profile
  config are never overwritten. Base ships no seed → no-op.
- non-root (`node`), tini, `DSH_TELEMETRY_DISABLED=1`, no browser/X11 stack

## Remote access

dsh's fence (docs/user-guide/public-deployments): pages and static assets
load for any authority, but **every /api call and WebSocket upgrade from
another authority is 403** unless named. Set:

- `DSH_TRUSTED_HOSTS` — comma-separated authorities (host or host:port;
  the entrypoint emits a repeatable `--trusted-host` per entry)
- `DSH_PUBLIC_URL` — advertised external root; the printed startup-token
  URL then points at the public address directly

The startup token itself **cannot be disabled** (it is the core browser-auth
mechanism — no off switch upstream), but the signed cookie persists 30 days
and survives restarts (secret in `$DSH_HOME/.credentials.yaml`, on the
volume), so only new devices need a fresh token URL from the logs.

## Usage

Standalone stock dsh web: `docker run -p 3080:3080 ghcr.io/soulwhisper/deepseek-harness`
(dsh listens on 3080 loopback; the relay publishes it on the container IP).
As a child base: `FROM deepseek-harness` with bake `--build-context`
(see `apps/dsh-web/docker-bake.hcl`).

## Child images

| Image | Profile plugins |
| ----- | --------------- |
| `dsh-web` | `dsh-ui-auth`, `dshmarket` (runtime keeps full install/update ability) |
| `dsh-tavern` | dsh-tavern core (custom recipe) + wrongbook + tavern-image + knock-knock + narrative-anchor + `dsh-ui-auth` (frozen: no runtime pnpm) |
