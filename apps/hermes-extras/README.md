## Hermes-agent plugin deps

Bundles the python packages that hermes-agent platform plugins need but the
upstream image venv does not ship — currently the Feishu platform deps
(`lark-oapi` + `requests-toolbelt`). Intended as a Kubernetes init-container:
runtime lazy-install is impossible in-cluster (read-only rootfs), and plugin
deps belong in GitOps content, not runtime fetches.

Versions are pinned in [`docker-bake.hcl`](./docker-bake.hcl) and tracked by
renovate (`datasource=pypi`). The closure is installed with `--no-deps` on
purpose: every transitive dep (requests, httpx, websockets, pycryptodome,
qrcode, pypng) already exists in the upstream venv, and a copy here would
shadow it via PYTHONPATH. When bumping `lark-oapi`, re-check the closure
against the target hermes-agent venv.

### How it lands

[`app/install.sh`](./app/install.sh) exact-mirrors `/pylib` into
`$PYLIB_DIR` (default `/opt/data/.local/lib/hermes-plugins`, on the agent's
PVC). Stale copies of previous pins are evicted first. The agent container
sets:

```yaml
env:
  - name: PYTHONPATH
    value: /opt/data/.local/lib/hermes-plugins
```

```yaml
spec:
  initContainers:
    - name: hermes-extras
      image: ghcr.io/soulwhisper/hermes-extras:latest
      command: ["/bin/sh", "/app/install.sh"]
      volumeMounts:
        - name: data
          mountPath: /opt/data
```

### Discovery

The default `CMD` lists the bundle:

```console
$ docker run --rm ghcr.io/soulwhisper/hermes-extras:latest
lark_oapi
lark_oapi-1.6.8.dist-info
requests_toolbelt
requests_toolbelt-1.0.0.dist-info
```
