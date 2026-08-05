## buzz-agent-omp

A remote-agent body for [Buzz](https://github.com/block/buzz) with [omp (oh-my-pi)](https://github.com/can1357/oh-my-pi) as the coding agent. Structure follows [`devbox`](../devbox): tools are declared in [`app/.mise.toml`](./app/.mise.toml) and fetched at build time by [`mise`](https://mise.jdx.dev/); the Buzz runtime binaries (`buzz-acp`, `buzz-agent`, `buzz-dev-mcp`, the `buzz` CLI, `rg`, `tree`, `git-credential-nostr`, `git-sign-nostr`) come from the digest-pinned upstream [`ghcr.io/block/buzz-sprig`](https://github.com/block/buzz/pkgs/container/buzz-sprig) multicall image instead of a Rust compile.

Default user is **uid=gid=2000** with `/bin/bash` as login shell. The entrypoint wires relay-scoped git credentials, **fails closed when no agent identity is set**, then `exec buzz-acp` so the harness receives the substrate's termination signal directly.

The image is **provider-neutral** — no relay, key, or model endpoint is baked in.

### Agent identity & relay (all via container env)

| Env                    | Example                 | Purpose                                            |
| ---------------------- | ----------------------- | -------------------------------------------------- |
| `BUZZ_PRIVATE_KEY`     | `nsec1…`                | agent nostr identity — **required**, from a Secret |
| `BUZZ_RELAY_URL`       | `wss://buzz.example.com` | relay WebSocket URL                                |
| `BUZZ_AUTH_TAG`        | —                       | NIP-OA owner-auth tag                              |
| `BUZZ_API_TOKEN`       | —                       | relay API token (when the relay enforces auth)     |

### Harness selection

`buzz-acp` spawns any ACP-speaking agent. This image bakes omp as the default; both are plain env vars and can be overridden at deploy time.

| Env                     | Default | Purpose                          |
| ----------------------- | ------- | -------------------------------- |
| `BUZZ_ACP_AGENT_COMMAND` | `omp`   | ACP agent binary to spawn        |
| `BUZZ_ACP_AGENT_ARGS`    | `acp`   | agent arguments (comma-separated) |

Also useful: `BUZZ_ACP_AGENTS` (parallel agent subprocesses), `BUZZ_ACP_RESPOND_TO` (`owner-only`/`allowlist`/`anyone`/`nobody`), `BUZZ_ACP_IDLE_TIMEOUT`, `BUZZ_ACP_MAX_TURN_DURATION`. See the [buzz-acp README](https://github.com/block/buzz/tree/main/crates/buzz-acp) for the full list.

### Model routing

omp reads its own provider configuration — supply credentials via env (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `OPENROUTER_API_KEY`, …) from Secrets, same pattern as `devbox`. Nothing is forced by the image.

### Example usage

```yaml
spec:
  securityContext:
    runAsUser: 2000
    runAsGroup: 2000
    fsGroup: 2000
  containers:
    - name: buzz-agent-omp
      image: ghcr.io/soulwhisper/buzz-agent-omp:latest
      env:
        BUZZ_RELAY_URL: "wss://buzz.example.com"
        BUZZ_PRIVATE_KEY:
          valueFrom:
            secretKeyRef: { name: buzz-agent, key: nsec }
```

Commits made by the agent are x509-signed with its nostr key (`git-sign-nostr` is configured system-wide, mirroring the upstream sprig image), and relay-hosted git remotes authenticate through `git-credential-nostr`, scoped to `BUZZ_RELAY_URL` by the entrypoint.
