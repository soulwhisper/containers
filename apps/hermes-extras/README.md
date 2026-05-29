## Hermes-agent Extra binaries

Bundle [hermes-agent](https://github.com/NousResearch/hermes-agent) extra binaries into a single image, intended for use as a Kubernetes init-container in network-restricted environments where the agent pod cannot reach upstream release hosts (GitHub, custom registries, etc).

Tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time by [`mise`](https://mise.jdx.dev/). All resulting binaries are flattened into `/data/` for predictable mounting.

### Example usage

```yaml
spec:
  initContainers:
    - name: hermes-extras
      image: ghcr.io/soulwhisper/hermes-extras:latest
      command:
        [
          "sh",
          "-c",
          "cp -v /data/* /opt/data/.local/bin/ && chmod +x /opt/data/.local/bin/*",
        ]
      volumeMounts:
        - name: extras
          mountPath: /opt/data/.local/bin
  containers:
    - name: hermes-agent
      image: ghcr.io/soulwhisper/hermes-suite:latest
      env:
        - name: PATH
          value: /opt/data/.local/bin:/opt/hermes/.venv/bin:/usr/local/bin:/usr/bin:/bin
      volumeMounts:
        - name: extras
          mountPath: /opt/data/.local/bin
  volumes:
    - name: extras
      emptyDir: {}
```

### Discovery

The default `CMD` prints a version table of every bundled binary:

```sh
docker run --rm ghcr.io/soulwhisper/hermes-extras:latest
```
