## dify-sandbox

A sandbox container for [Dify](https://dify.ai/) workflow execution with **node** and **officecli** pre-installed via [mise](https://mise.jdx.dev/). Tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time.

Default user is **uid=gid=2000** with `/bin/bash` as login shell. Idles on `sleep infinity` so you can `kubectl exec -it -- bash`; pass argv to run a command directly instead.

### Tools

| Tool | Source | Manager |
|------|--------|---------|
| node | mise built-in | mise |
| officecli | npm `@officecli/officecli` | mise (npm backend) |

### Example usage

```yaml
spec:
  securityContext:
    runAsUser: 2000
    runAsGroup: 2000
    fsGroup: 2000
  containers:
    - name: sandbox
      image: ghcr.io/soulwhisper/dify-sandbox:latest
```

```sh
kubectl exec -it deploy/dify-sandbox -- bash
kubectl exec -it deploy/dify-sandbox -- officecli --version
```
