## Hindsight — local models baked in

A derived [Hindsight](https://github.com/vectorize-io/hindsight) image that
bundles a local embedder and reranker so the container needs no HuggingFace
access at runtime and no model-cache PVC. CPU-only torch keeps the layer set
lean; the registry handles per-node caching of the weights.

| Component | Model                     | Notes                          |
| --------- | ------------------------- | ------------------------------ |
| embedder  | `BAAI/bge-m3`             | ~2.3GB, 1024-dim, multilingual |
| reranker  | `BAAI/bge-reranker-v2-m3` | ~1.1GB, cross-encoder          |

Built from `ghcr.io/vectorize-io/hindsight:<version>-slim`. Service ports are
inherited from upstream: `8888` (API) and `9999` (control-plane UI).

### Versioning

- Image tag follows the upstream **slim** release with the `-slim` suffix
  stripped (e.g. upstream `0.8.3-slim` → `ghcr.io/soulwhisper/hindsight:0.8.3`).
- Renovate tracks `VERSION` in `docker-bake.hcl` via the `docker` datasource
  with `versioning=docker`, so updates stay on the `*-slim` line.
- The bundled model set is pinned in the `Dockerfile` (`EMBEDDER` / `RERANKER`
  build args). Set `EMBEDDER_REV` / `RERANKER_REV` to a HF commit sha for fully
  reproducible weight bake-in.

### Baked-in defaults

The image ships local-provider defaults so it works out of the box:

```
HINDSIGHT_API_EMBEDDINGS_PROVIDER=local
HINDSIGHT_API_EMBEDDINGS_LOCAL_MODEL=BAAI/bge-m3
HINDSIGHT_API_RERANKER_PROVIDER=local
HINDSIGHT_API_RERANKER_LOCAL_MODEL=BAAI/bge-reranker-v2-m3
HF_HOME=/app/models
OMP_NUM_THREADS=4
MKL_NUM_THREADS=4
```

These are defaults only — any value set in the pod spec / HelmRelease env
overrides them, so runtime config stays in GitOps. Keep `OMP_NUM_THREADS` /
`MKL_NUM_THREADS` aligned with `resources.limits.cpu` to avoid throttling and
noisy-neighbour behaviour on shared nodes.

### Immutable / read-only rootfs

Weights live under `/app/models` (world-readable) and are only ever read at
runtime, so the image is compatible with `readOnlyRootFilesystem: true`
provided the embedded database path (`/home/hindsight/.pg0`) is a writable
mount.

### Example usage

#### Docker Compose

```yaml
services:
  hindsight:
    image: ghcr.io/soulwhisper/hindsight:0.8.3
    container_name: hindsight
    restart: unless-stopped
    ports:
      - "8888:8888" # api
      - "9999:9999" # control plane
    volumes:
      - hindsight-data:/home/hindsight/.pg0
    environment:
      - HINDSIGHT_API_LLM_PROVIDER=openai
      - HINDSIGHT_API_LLM_API_KEY=sk-xxx

volumes:
  hindsight-data:
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
        - name: hindsight
          image: ghcr.io/soulwhisper/hindsight:0.8.3
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          resources:
            limits:
              cpu: "4"
          ports:
            - { name: api, containerPort: 8888 }
            - { name: control-plane, containerPort: 9999 }
          volumeMounts:
            - { name: data, mountPath: /home/hindsight/.pg0 }
```
