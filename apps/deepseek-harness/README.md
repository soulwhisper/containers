# deepseek-harness

DeepSeek Harness ([@deepseek-ai/dsh](https://www.npmjs.com/package/@deepseek-ai/dsh)) appliance image with two switchable profiles. Built from npm + github sources only — no third-party packaged base, no upstream fork.

| Profile | Purpose | Bundled plugins | Status |
| ------- | ------- | --------------- | ------ |
| `web` (default) | agent workbench (nix-ops) | `dsh-ui-auth` (login gateway), `dsh-pocket-relay` (self-hosted phone relay), `dshmarket` (plugin discovery) | ✅ verified end-to-end (boot, relay, login) |
| `tavern` | RP chat (cluster appliance) | [`flizzywine/dsh-tavern`](https://github.com/flizzywine/dsh-tavern) (SillyTavern-compatible runtime), `dsh-ui-auth`, `dsh-pocket` (mobile access) | ⚠️ baked but blocked upstream — see below |

### Known issue: tavern plugin import fails on dsh 0.2.1-rc.2

`dsh-tavern` (v3.0 and `#main`, also `#v2.7.0-dsh0.2.1-test.2`) fails to
import its plugin entry in the CLI profile flow:

```
dsh-tavern (…/tavern-plugin/lib/index.js): failed to import:
TypeError: createRequire.resolve.paths is not a function
    at ResolutionRouter.routeScoped (@deepseek-ai/dsh-app-boot)
```

Reproduced in a **vanilla** `node:22-trixie` container with the documented
official flow (no custom image content), so it is an upstream loader issue,
not an image bug — tracked at
[flizzywine/dsh-tavern#300](https://github.com/flizzywine/dsh-tavern/issues/300).
The bundle stays baked on purpose: once a fixed dsh/tavern release lands,
bumping `DSH_VERSION` / `DSH_TAVERN_REF` is the only change needed.
Non-tavern profiles are unaffected (web verified 302 through the relay).

Select the profile with `DSH_PROFILE=web|tavern` or as `$1`. All plugins are
baked at build time into `/opt/dsh-seed` (see Dockerfile); runtime egress is
not required for plugin operation.

## Runtime contract

- dsh listens on `127.0.0.1:3080`; `app/dsh-web-relay.cjs` binds the
  container's concrete IPv4 (multi-NIC ambiguity requires
  `DSH_WEB_RELAY_INTERFACE`) so a k8s Service or LAN can reach it.
- State lives in `$DSH_HOME` (`/home/node/.dsh`, mount a volume there).
  User data = `profile-data/`; plugin code = `profiles/`.
- Seeding (entrypoint, idempotent):
  - fresh `$DSH_HOME` → whole seed copied (profiles config included);
  - image bump (`BUNDLED_VERSIONS` changed) → only plugin **code**
    (`profiles/*/node_modules`) is re-synced; `profile-data/` and profile
    config files are never overwritten. Runtime `dsh plugin add` additions
    are *not* preserved across image updates — bundle changes go through
    the image.
- First login: configure the model provider in the tavern/dsh settings UI
  (persisted on the volume). `dsh-ui-auth` provides the login gateway;
  behind a reverse proxy set `DSH_AUTH_TRUST_PROXY=1`.
- `DSH_TELEMETRY_DISABLED=1` is set in the image.

## Versioning & updates

| Pin | Where | datasource |
| --- | ----- | ---------- |
| `@deepseek-ai/dsh` | Dockerfile ARG + bake tag | npm |
| `pnpm` | Dockerfile ARG | npm |
| `flizzywine/dsh-tavern` | Dockerfile ARG `DSH_TAVERN_REF` | github-tags |
| `dsh-ui-auth` / `dsh-pocket` / `dsh-pocket-relay` / `dshmarket` | Dockerfile ARGs | npm |
| `node:24-trixie` | Dockerfile `NODE_IMAGE` ARG | docker |

dsh does not self-update inside the container; plugins do not self-update
either (upstream semantics: upgrade = re-run `dsh plugin add`). The image is
the update vehicle: bump a pin, rebuild, roll. `engines.dsh` (plugins) and
the tavern startup check gate misalignment at boot instead of failing
mid-session.

## Single image vs two sub-images (evaluation)

Chosen: **one image, runtime profile switch**. Reasoning:

- The profile state (`profiles/`, plugin code) is seeded per-profile at
  runtime anyway — a per-profile image would only change which subset of
  `node_modules` ships, i.e. save the *other* profile's plugin bytes.
- Measured on this build: image **2062 MiB** total; dsh core
  (`@deepseek-ai/dsh` tree) **840 MiB**; baked plugin seed **136 MiB** of
  which tavern-specific ≈ **116 MiB**, web-specific ≈ **21 MiB**. Splitting
  saves each role ≈ 1–6 % — far below the 30 % trip line. The elephant is
  the dsh core, shared by both roles.
- One artifact = one pin set, one scan/verify path, identical auth/relay
  behavior in both roles; role separation is enforced at deploy time
  (different args, CNP, PVCs).

Split to two sub-images (`deepseek-harness-web`, `deepseek-harness-tavern`)
becomes the right call when **any** of these trips (revisit per release):

1. role-specific plugin payload exceeds ~30 % of image size (measure with
   `docker image inspect` — print current value in the release notes);
2. a heavyweight optional appears (browser/X11 stack — that gets its own
   `desktop` build *target*, never into the default image);
3. the two roles need different dsh versions (tavern pins DSH compat
   aggressively; the workbench may want to track faster).

Implementation note if splitting: keep this Dockerfile, move profile bake
into `AS plugins-web` / `AS plugins-tavern` stages and add bake targets
`web`/`tavern`; the entrypoint contract is unchanged.

## Local build / test

```bash
docker buildx bake -f docker-bake.hcl --set image-local.args.GIT_SHA=$(git rev-parse HEAD) ./apps/deepseek-harness
APP=deepseek-harness TEST_IMAGE=deepseek-harness:0.2.1-rc.2 go test -v ./apps/deepseek-harness/...
```
