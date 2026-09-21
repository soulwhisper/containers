## God's Eye View

A containerized [God's Eye View](https://github.com/bilawalsidhu/gods-eye-view)
— live open-source spatial intelligence (aircraft, vessels, satellites, CCTV,
traffic) on a photorealistic 3D globe.

Upstream ships no image and no standalone production server: the supported
runtime is a Vite dev server that serves the client and brokers every
third-party API same-origin under `/api/*`, exactly like upstream's own
headless launcher (`scripts/pinokio-start.mjs`). This image mirrors that
launcher: full upstream checkout at the pinned ref plus `node_modules`, keys
via env. `HOST=0.0.0.0` is set so upstream's vite config enables
`allowedHosts: true` — required when a reverse proxy forwards a public Host
header.

Service port: `4173` (web UI + `/api/*` provider proxies).

### Versioning

- `VERSION` in `docker-bake.hcl` is the upstream git ref (`main` while
  upstream has no release tags); the image also gets `sha-<containers-sha>`
  and `latest` tags from CI.
- Puppeteer (QA-only dep) Chromium download is skipped at build time.

### Runtime configuration

All keys are optional — without any key the app boots on keyless Esri World
Imagery. Keys are consumed from the environment at server start (vite `define`
injects the two client-exposed ones into the browser bundle):

| Env var                    | Scope         | Purpose                              |
| -------------------------- | ------------- | ------------------------------------ |
| `GOOGLE_MAPS_API_KEY`      | client-exposed | Google 3D Tiles + GEV place search  |
| `CESIUM_ION_TOKEN`         | client-exposed | ion tiles/terrain fallback          |
| `GOOGLE_MAPS_SERVER_API_KEY` | server-only  | Places + Street View fallback       |
| `OPENAI_API_KEY`           | server-only   | Realtime voice control               |
| `OPENSKY_CLIENT_ID/SECRET` | server-only   | OpenSky aircraft (OAuth mode)        |
| `AISSTREAM_API_KEY`        | server-only   | Live AIS vessel stream (websocket)   |
| `FIRMS_MAP_KEY`            | server-only   | NASA active fires                    |
| `TOMTOM_API_KEY`           | server-only   | Live traffic flow                    |
| `LL2_API_TOKEN`            | server-only   | Launch Library 2 allowance           |

Built-in rate limits default to on (`GEV_RATELIMIT_OPENAI_PER_MIN=30`,
`GEV_RATELIMIT_GOOGLE_PER_MIN=60`) per upstream guidance for non-localhost
serving.

### Read-only rootfs

The app writes only vite's optimize-deps cache (`/app/node_modules/.vite`)
and the voice debug log (`/app/.gev-logs`). Mount emptyDir volumes on both to
run with `readOnlyRootFilesystem: true`.

### Example usage

#### Docker Compose

```yaml
services:
  gods-eye-view:
    image: ghcr.io/soulwhisper/gods-eye-view:latest
    container_name: gods-eye-view
    restart: unless-stopped
    ports:
      - "4173:4173"
    environment:
      - GOOGLE_MAPS_API_KEY=xxx
```

#### Kubernetes (non-root, read-only rootfs)

```yaml
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - name: app
          image: ghcr.io/soulwhisper/gods-eye-view:latest
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          ports:
            - { name: http, containerPort: 4173 }
          volumeMounts:
            - { name: vite-cache, mountPath: /app/node_modules/.vite }
            - { name: gev-logs, mountPath: /app/.gev-logs }
            - { name: tmp, mountPath: /tmp }
      volumes:
        - { name: vite-cache, emptyDir: {} }
        - { name: gev-logs, emptyDir: {} }
        - { name: tmp, emptyDir: {} }
```
