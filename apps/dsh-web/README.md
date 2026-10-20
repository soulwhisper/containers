# dsh-web

Agent workbench appliance (nix-ops) on the `deepseek-harness` basement.

| Profile | Bundled plugins (build-time) | Runtime ability |
| ------- | ---------------------------- | --------------- |
| `web` (default) | `dsh-ui-auth` (login gateway: sessions survive restarts, reverse-proxy ready, 2FA) + `dshmarket` (in-UI plugin discovery) | **full**: pnpm stays in the image; runtime `dsh plugin add` works (container-test contract). Add what you need via dshmarket or CLI |

## Runtime contract

Same as the base (see `apps/deepseek-harness/README.md`): relay publishes
dsh's loopback web on the container IP (`DSH_PORT`, upstream
`DSH_UPSTREAM_PORT`); `$DSH_HOME` on a volume keeps settings, sessions and
runtime-installed plugins; first boot seeds the baked profile, image bumps
re-sync plugin code only (`profile-data/` and runtime installs are
untouched — plugins installed at runtime persist across image updates
because only the baked `node_modules` set is re-synced... note: runtime adds
live in the same `profiles/web/node_modules`, and the seed sync uses
`rsync --delete` **only on that directory** — a runtime-installed plugin
will be removed by an image bump unless also bundled. Keep wanted plugins in
this image (PR a pin) if they must survive).

Behind a reverse proxy set `DSH_AUTH_TRUST_PROXY=1`.

## Versioning

| Pin | datasource |
| --- | ---------- |
| base image | `../deepseek-harness` (bake target `dsh-base`) |
| `dsh-ui-auth` | npm |
| `dshmarket` | npm |

## Build

```bash
docker buildx bake -f docker-bake.hcl ./apps/dsh-web   # builds ../deepseek-harness first
```
