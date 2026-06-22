## devbox

A long-running, cluster-internal dev container for **git + [Claude Code](https://github.com/anthropics/claude-code)**. Tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time by [`mise`](https://mise.jdx.dev/); `git`/`bash`/`openssh` come from the base image.

Default user is **uid=gid=2000** with `/bin/bash` as login shell. Idles on `sleep infinity` so you can `kubectl exec -it -- bash`; pass argv (`claude`, `git ...`, a CI step) to run it directly instead.

The image is **provider-neutral** — no model endpoint is baked in, so it never forces a backend.

### Model routing (all via container env)

`ANTHROPIC_BASE_URL` must reach an Anthropic Messages-compatible route on agentgateway; `ANTHROPIC_AUTH_TOKEN` is the gateway bearer token. Model ids are whatever the gateway maps.

| Env                                                                 | Example                              | Purpose                            |
| ------------------------------------------------------------------- | ------------------------------------ | ---------------------------------- |
| `ANTHROPIC_BASE_URL`                                                | `https://api.deepseek.com/anthropic` | gateway Anthropic route            |
| `ANTHROPIC_AUTH_TOKEN` / `ANTHROPIC_API_KEY`                        | —                                    | gateway credential (from a Secret) |
| `ANTHROPIC_MODEL`                                                   | `deepseek-v4-pro[1m]`                | main model                         |
| `ANTHROPIC_DEFAULT_OPUS_MODEL` / `…_SONNET_MODEL` / `…_HAIKU_MODEL` | —                                    | model-tier aliases                 |
| `CLAUDE_CODE_SUBAGENT_MODEL`                                        | `deepseek-v4-flash`                  | subagent model                     |

### Example usage

```yaml
spec:
  securityContext:
    runAsUser: 2000
    runAsGroup: 2000
    fsGroup: 2000
  containers:
    - name: devbox
      image: ghcr.io/soulwhisper/devbox:latest
      env:
        ANTHROPIC_BASE_URL: "https://api.deepseek.com/anthropic"
        ANTHROPIC_MODEL: "deepseek-v4-pro[1m]"
        ANTHROPIC_DEFAULT_OPUS_MODEL: "deepseek-v4-pro[1m]"
        ANTHROPIC_DEFAULT_SONNET_MODEL: "deepseek-v4-pro"
        ANTHROPIC_DEFAULT_HAIKU_MODEL: "deepseek-v4-flash"
        CLAUDE_CODE_SUBAGENT_MODEL: "deepseek-v4-flash"
        ANTHROPIC_AUTH_TOKEN:
          valueFrom:
            secretKeyRef: { name: llm-api, key: deepseek }
```

```sh
kubectl exec -it deploy/devbox -- bash
```
