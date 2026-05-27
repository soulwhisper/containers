## Hermes-agent Plugins

Bundle [hermes-agent](https://github.com/NousResearch/hermes-agent) plugin binaries into a single image, intended for use as a Kubernetes init-container in network-restricted environments where the agent pod cannot reach plugins' upstream release hosts (GitHub, custom registries, etc).

Tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time by [`mise`](https://mise.jdx.dev/). All resulting binaries are flattened into `/plugins/` for predictable mounting.

### Example usage

#### Kubernetes init-container

```yaml
spec:
  initContainers:
    - name: hermes-plugins
      image: ghcr.io/soulwhisper/hermes-plugins:latest
      command: ["sh", "-c", "cp -v /plugins/* /target/ && chmod +x /target/*"]
      volumeMounts:
        - name: plugins
          mountPath: /target
  containers:
    - name: hermes-agent
      image: ghcr.io/soulwhisper/hermes-suite:latest
      env:
        - name: PATH
          value: /opt/data/plugins:/opt/hermes/.venv/bin:/usr/local/bin:/usr/bin:/bin
      volumeMounts:
        - name: plugins
          mountPath: /opt/data/plugins
  volumes:
    - name: plugins
      emptyDir: {}
```

#### Discovery

The default `CMD` prints a version table of every bundled binary:

```sh
docker run --rm ghcr.io/soulwhisper/hermes-plugins:latest
```

### Adding a new plugin

1. Add the tool to `app/.mise.toml` — `mise` supports many backends (`github:owner/repo`, `aqua:`, `ubi:`, etc.); see the [mise registry](https://mise.jdx.dev/registry.html).
2. Run `just local-build hermes-plugins` and confirm the binary lands in `/plugins/`.
3. Add an existence test in `container_test.go`.
4. If the new tool should drive the image tag, replace the `rtk-ai/rtk` reference in `docker-bake.hcl` and the renovate grouping rule accordingly. Otherwise the existing rtk-based tag remains the image's primary version label.
