# dsh-tavern

RP chat appliance (cluster) on the `deepseek-harness` basement:
[flizzywine/dsh-tavern](https://github.com/flizzywine/dsh-tavern) profile +
ecosystem plugins, **frozen** (no runtime plugin changes; updates only via
image rebuilds).

## Bundled content

| Piece | Pin | Notes |
| ----- | --- | ----- |
| dsh-tavern core | `v3.0` | **custom recipe**, not install.sh: clone → `pnpm install --frozen-lockfile` (workspace deps — the step plain `dsh plugin add` misses, see dsh-tavern#300) → generated profile manifest + pruned workspace yaml → `pnpm install` in profile. Recipe is build-asserted at every step |
| `dsh-tavern-wrongbook` | `v1.33.0` | per-card debugging ledger (activation lines asserted in boot smoke) |
| `dsh-tavern-image` | `v2.9.0` | scene illustration; backend (ComfyUI/online keys) configured in UI, stored on the volume |
| `dsh-tavern-knock-knock` | `v0.7.1` | long-campaign reminders; installs as plain dependency (upstream warning expected, not an error) |
| `narrative-anchor` | `#2ee3da04` | no release tags on the current line — pinned commit, bump manually. First boot logs a benign ENOENT for a desktop-layout audit path (expected) |
| auth | — | **none in-app** — auth is the authentik forward-auth portal at the gateway (united auth). Fallback: the dsh startup-token URL in pod logs; the token rotates per process start, so a pod restart re-locks the UI until the new token URL is used once |

| host UX patches | from the tavern source | `theme-font-limit` (text max 22→32) + `chat-answer-step` (edited rounds stay expanded) applied to this image's dsh tree only; idempotent and skip-safe on dsh bumps (build prints their result) |

## Freeze contract

`pnpm`/`pnpx` binaries are removed from the runtime image. Runtime
`dsh plugin add` (and the in-UI updater) fail fast — that is the stability
design: the only change channel is rebuilding this image. Persistent state
(cards, sessions, settings, plugin config) lives under `$DSH_HOME`
(`profile-data/` + profile config) on the volume and is never overwritten by
the seed sync.

## Seed / update semantics

Shared with the base entrypoint: fresh `$DSH_HOME` gets the whole baked
profile (`/opt/dsh-seed`); an image bump (`BUNDLED_VERSIONS` diff) re-syncs
plugin **code** only. The tavern source tree stays at
`/opt/dsh-tavern/apps/dsh-tavern` — the profile links into it
(`link:` deps), so it ships in the image read-only.

## Versioning

| Pin | datasource |
| --- | ---------- |
| base image | `../deepseek-harness` (bake target `dsh-base`) |
| `DSH_TAVERN_REF` | github-tags `flizzywine/dsh-tavern` |
| `DSH_WRONGBOOK_REF` / `DSH_TAVERN_IMAGE_REF` / `DSH_KNOCK_KNOCK_REF` | github-tags |
| `DSH_NARRATIVE_ANCHOR_SHA` | manual (see table) |
| `DSH_UI_AUTH_VERSION` | npm |

When bumping `DSH_TAVERN_REF`, re-check the source `patches/` and
`pnpm-workspace.yaml`: the profile workspace is **pruned to the installed
dependency set** (source yaml references pocket/web-mobile patches that fail
`ERR_PNPM_UNUSED_PATCH` under pnpm 11). Pruning rule lives in
`app/profile-pnpm-workspace.yaml`.

## Build

```bash
docker buildx bake -f docker-bake.hcl ./apps/dsh-tavern   # builds ../deepseek-harness first
```
