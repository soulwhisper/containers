## hermes-agent — upstream image plus baked extras

The upstream [`nousresearch/hermes-agent`](https://github.com/NousResearch/hermes-agent)
image, augmented at build time:

- **supervisord** running `gateway` + `dashboard` in one container (upstream's
  own s6-overlay entrypoint needs root + a writable `/run`, incompatible with
  `runAsNonRoot` / dropped caps / read-only rootfs)
- **caveman** skill pack (`/bundled/skills`, exact-mirrored onto the data
  volume by the entrypoint — hermes loads skills natively, enabled on load)
- **rtk** — shell-output compression proxy for agent commands
  (kubectl/git/docker/…). Fully initialized at build: binary at
  `/bundled/bin` (on `PATH`, sha256-verified) and the hermes plugin adapter
  at `/bundled/plugins/rtk-rewrite` (mirrored to `~/.hermes/plugins/` by the
  entrypoint). The only runtime contract is the allow-list —
  `plugins.enabled: [rtk-rewrite]` in the GitOps-managed `config.yaml` (a
  runtime `rtk init` would be reverted by the next config seed anyway)
- **feishu plugin deps** (`lark-oapi`, `qrcode`, …) installed into the project
  venv using upstream's own pins (`pyproject.toml` →
  `optional-dependencies.feishu`) — runtime lazy-install is impossible on a
  read-only rootfs

| Service     | Port | Notes                                            |
| ----------- | ---- | ------------------------------------------------ |
| `gateway`   | 8642 | agent API / feishu / cron / webhooks             |
| `dashboard` | 9119 | built-in web UI; auth gate requires a provider   |

Services are selected via `HERMES_SERVICES` (default `gateway,dashboard`); a
single selected service execs directly as PID 1, skipping supervisord.

### Versioning

Only **hermes-agent** is version-tracked:

- The image tag follows upstream's semver GitHub releases (e.g. `0.21.7`).
- The base image is `nousresearch/hermes-agent:stable@sha256:…` — Docker Hub
  publishes no plain `vX.Y.Z` tag (only `rc.N-vX.Y.Z` candidates and floating
  aliases), so the digest is pinned and floated by renovate. A renovate group
  with `minimumGroupSize: 2` ships the tag bump and the digest bump in one PR.
- Upstream's legacy `v2026.x` CalVer GitHub releases still exist and sort
  newer under semver; a packageRule blocks them (`!/^\d{4}\./`).

Extras (**rtk**, **caveman**) are deliberately **untracked**: they resolve to
latest at build time and ride along on every hermes-agent rebuild without
triggering builds themselves. Resolved versions are recorded in
`/bundled/VERSIONS`. The CI gate is behavioral: the test suite runs `rtk
--version` and imports the feishu deps inside the built image, so a bad
floating resolution fails before publish.

### Non-root by default

The image's default user is `hermes` (UID 10000), compatible with
`runAsNonRoot: true` + `readOnlyRootFilesystem: true` (provide writable mounts
at `/opt/data`, `/workspace`, `/tmp`). Starting as root is still supported for
docker/compose: the entrypoint remaps via `HERMES_UID`/`HERMES_GID`, chowns,
then drops via `gosu`.

supervisord state lives under `$HERMES_HOME/.supervisor/`:

```sh
docker exec hermes-agent supervisorctl -c /app/supervisord.conf status
```

### Example usage

```yaml
services:
  hermes-agent:
    image: ghcr.io/soulwhisper/hermes-agent:0.21.7
    restart: unless-stopped
    ports:
      - "8642:8642" # gateway
      - "9119:9119" # dashboard
    volumes:
      - ./hermes:/opt/data
      - ./workspace:/workspace
```

On first run the entrypoint seeds `/opt/data/.env`, `config.yaml` and
`SOUL.md` from the bundled examples. The dashboard refuses to bind
non-loopback without an auth provider — set `dashboard.basic_auth` in
`config.yaml` (or run behind an authenticating gateway, as home-ops does).
