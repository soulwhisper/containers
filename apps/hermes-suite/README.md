## Hermes Suite — all-in-one image

Combines three services into a single image. Services are selected at startup
via `HERMES_SERVICES`; when only one is selected the container `exec`s the
binary directly as PID 1 (no supervisord). Otherwise `supervisord` runs the
selected combination.

| Service     | Port | Source                                                 |
| ----------- | ---- | ------------------------------------------------------ |
| `gateway`   | 8642 | hermes-agent CLI / Telegram / cron / webhook           |
| `dashboard` | 9119 | built-in monitoring dashboard                          |
| `webui`     | 8787 | chat [webui](https://github.com/nesquena/hermes-webui) |

### Service selection

| `HERMES_SERVICES`         | Result                                |
| ------------------------- | ------------------------------------- |
| `gateway`                 | gateway only, PID 1 = hermes binary   |
| `webui`                   | webui only, PID 1 = python server     |
| `dashboard`               | dashboard only, PID 1 = hermes binary |
| `gateway,webui` (default) | supervisord runs both                 |
| `gateway,webui,dashboard` | supervisord runs all three            |

Unknown service names cause the entrypoint to exit non-zero with a clear error.

### Non-root by default

The image's default user is `hermes` (non-root, UID 1000). Compatible with
k8s `securityContext.runAsNonRoot: true` and `readOnlyRootFilesystem: true`
out of the box (provided `/opt/data` is a writable mount).

- **k8s**: set `securityContext.fsGroup: 10000` so the PVC at `/opt/data` is
  group-writable.
- **docker / compose**: pre-create the host directory with `chown -R 10000:10000 ./hermes ./workspace`,
  or start the container once as root (`--user 0`) — the entrypoint will fix
  ownership automatically and then drop to `hermes` via `gosu`.

Root mode is preserved for parity with the previous image: if the container
starts as root, the entrypoint optionally remaps the bundled `hermes` user
via `HERMES_UID` / `HERMES_GID`, chowns `/opt/data` and `/workspace`, then
`exec gosu hermes` to re-run itself as `hermes`.

supervisord runtime state (sock / pid / log) lives under
`$HERMES_HOME/.supervisor/`. Inspect with:

```sh
docker exec hermes-suite supervisorctl -c /app/supervisord.conf status
```

### Versioning

- Image tag follows the upstream **agent** CalVer (e.g. `2026.5.7`).
- Bundled **webui** version is recorded only in OCI labels:
  - `hermes-suite.agent-version`
  - `hermes-suite.webui-version`
- A webui-only bump rewrites `docker-bake.hcl`, which triggers the release
  workflow and republishes the same agent-CalVer tag with new webui content
  (new digest, updated label). Consumers who need exact pinning should pin
  by digest.

### Example usage

#### Docker Compose

```yaml
services:
  hermes-suite:
    image: ghcr.io/soulwhisper/hermes-suite:2026.5.7
    container_name: hermes-suite
    restart: unless-stopped
    ports:
      - "8642:8642" # gateway
      - "8787:8787" # webui
    volumes:
      - ./hermes:/opt/data
      - ./workspace:/workspace
      - /etc/localtime:/etc/localtime:ro
    environment:
      - HERMES_SERVICES=gateway,webui
```

On first run the entrypoint seeds `/opt/data/.env`, `/opt/data/config.yaml`
and `/opt/data/SOUL.md` from the bundled examples. Edit `/opt/data/.env` to
add API keys.

#### Kubernetes (non-root, dashboard enabled)

```yaml
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 10000 # hermes
        runAsGroup: 10000
        fsGroup: 10000
      containers:
        - name: hermes-suite
          image: ghcr.io/soulwhisper/hermes-suite:2026.5.7
          env:
            - name: HERMES_SERVICES
              value: gateway,webui,dashboard
          ports:
            - { name: gateway, containerPort: 8642 }
            - { name: webui, containerPort: 8787 }
            - { name: dashboard, containerPort: 9119 }
          volumeMounts:
            - { name: data, mountPath: /opt/data }
            - { name: workspace, mountPath: /workspace }
```

#### Single-service deployment

Useful when running each service as an independent workload sharing one PVC.
The entrypoint skips supervisord and execs the binary directly:

```yaml
env:
  - name: HERMES_SERVICES
    value: webui
```
