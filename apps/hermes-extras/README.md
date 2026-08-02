## Hermes-agent Extra binaries

Bundle [hermes-agent](https://github.com/NousResearch/hermes-agent) extra binaries into a single image, intended for use as a Kubernetes init-container in network-restricted environments where the agent pod cannot reach upstream release hosts (GitHub, custom registries, etc).

Tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time by [`mise`](https://mise.jdx.dev/). All resulting binaries are flattened into `/data/` for predictable mounting.

Also bundles the [caveman](https://github.com/JuliusBrussee/caveman) skill pack (pinned tarball, version managed via `CAVEMAN_VERSION` in [`docker-bake.hcl`](./docker-bake.hcl)) under `/skills/` — prompt skills, not binaries, so mise cannot fetch them.

### Example usage

The baked [`app/install.sh`](./app/install.sh) mirrors binaries into `$BIN_DIR` (default `/opt/data/.local/bin`) and the caveman skills into `$SKILLS_DIR` (default `/opt/data/skills/productivity`, i.e. `$HERMES_HOME/skills/productivity`). Idempotent; skills are exact-synced, foreign files in `BIN_DIR` are left alone.

```yaml
spec:
  initContainers:
    - name: hermes-extras
      image: ghcr.io/soulwhisper/hermes-extras:latest
      command: ["/bin/sh", "/app/install.sh"]
      volumeMounts:
        - name: extras
          mountPath: /opt/data
  containers:
    - name: hermes-agent
      image: ghcr.io/soulwhisper/hermes-suite:latest
      env:
        - name: PATH
          value: /opt/data/.local/bin:/opt/hermes/.venv/bin:/usr/local/bin:/usr/bin:/bin
      volumeMounts:
        - name: extras
          mountPath: /opt/data
  volumes:
    - name: extras
      persistentVolumeClaim:
        claimName: hermes-agent
```

### Discovery

The default `CMD` prints a version table of every bundled binary:

```sh
docker run --rm ghcr.io/soulwhisper/hermes-extras:latest
```
