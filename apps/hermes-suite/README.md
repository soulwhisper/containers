## Hermes Suite — all-in-one image

Combines three services into a single container, supervised by `supervisord`:

| Service            | Port | Source                                                 |
| ------------------ | ---- | ------------------------------------------------------ |
| `hermes-gateway`   | 8642 | hermes-agent CLI/Telegram/cron/webhook                 |
| `hermes-dashboard` | 9119 | built-in monitoring dashboard (part of hermes-agent)   |
| `hermes-webui`     | 8787 | chat [webui](https://github.com/nesquena/hermes-webui) |

This packaging avoids the Podman v3.4.4 limitation around sharing UID/GID
namespaces between multiple sibling containers — everything runs as a single
user inside one image.

### Versioning

Two upstream components are tracked independently and reflected in the image
tag as `<agent>-<webui>` (the leading `v` is stripped from each), for example
`2026.5.7-0.51.44`. Both are managed by Renovate via `docker-bake.hcl`:

- `AGENT_VERSION` — `docker.io/nousresearch/hermes-agent` tag
- `WEBUI_VERSION` — `github.com/nesquena/hermes-webui` release tag

### Example usage

```yaml
services:
  hermes-suite:
    image: ghcr.io/soulwhisper/hermes-suite:latest
    container_name: hermes-suite
    restart: unless-stopped
    ports:
      - "8642:8642" # gateway
      - "8787:8787" # webui
      - "9119:9119" # dashboard
    volumes:
      - ./hermes:/opt/data
      - ./workspace:/workspace
      - /etc/localtime:/etc/localtime:ro
    environment:
      # Optional: remap the in-container hermes user to match host ownership
      - HERMES_UID=1000
      - HERMES_GID=1000
```

On first run the entrypoint seeds `/opt/data/.env`, `/opt/data/config.yaml`, and
`/opt/data/SOUL.md` from the bundled examples. Edit `/opt/data/.env` to add API keys.
