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
