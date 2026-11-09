## covel

Web and server image for [ackness/covel](https://github.com/ackness/covel), an agentic AI RPG framework. The build uses the upstream release's own `docker/Dockerfile`; this repository tracks the release tag and adds the standard OCI metadata and GHCR publishing targets.

Published platforms: `linux/amd64` and `linux/arm64`. Tags include `latest`, the upstream semantic version, and `sha-<commit>`.

### Usage

```sh
export COVEL_MEDIA_TOKEN_SECRET="$(openssl rand -hex 32)"
docker volume create covel-data
docker run -d --name covel \
  --restart unless-stopped \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid \
  -p 127.0.0.1:3001:3001 \
  -e COVEL_BIND_HOST=0.0.0.0 \
  -e COVEL_HOME=/home/node/.covel \
  -e COVEL_MEDIA_TOKEN_SECRET \
  -e STORE_BACKEND=sqlite \
  -e SQLITE_PATH=/home/node/.covel/data/covel.db \
  -e MEDIA_ROOT=/home/node/.covel/data/media \
  -v covel-data:/home/node/.covel \
  ghcr.io/soulwhisper/covel:latest
```

Open `http://127.0.0.1:3001`. The container runs as the non-root `node` user. `/home/node/.covel` stores the SQLite database, media, installed worlds, and plugins; back up this volume before upgrades.

The default `self` deployment tier is a trusted personal instance, not a multi-tenant service. Use an authenticated TLS reverse proxy and configure `COVEL_ALLOWED_HOSTS` before exposing it to other machines. Model providers can be configured in the browser; an optional `llm.toml` can be mounted read-only at `/app/llm.toml`.

Covel may reject databases created by older releases instead of migrating them. Read the [upstream release notes](https://github.com/ackness/covel/releases) before changing image versions.
