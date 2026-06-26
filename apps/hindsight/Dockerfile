# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

# Hindsight with non-default local embedder/reranker baked in.
#   embedder : BAAI/bge-m3              (~2.3GB weights, 1024-dim, multilingual)
#   reranker : BAAI/bge-reranker-v2-m3  (~1.1GB weights, cross-encoder)
# CPU-only torch keeps the image lean; baked-in weights remove the runtime
# HuggingFace dependency so the registry handles per-node caching and you need
# no model-cache PVC.
#
# Service ports (inherited from upstream):
#   8888 — hindsight-api           (memory API)
#   9999 — hindsight-control-plane (browser UI)

# Upstream slim standalone image. ghcr tags are bare-version with a `-slim`
# suffix (e.g. 0.8.3-slim); GitHub *releases* are v-prefixed. The bake file
# tracks/overrides VERSION; renovate's dockerfile manager pins the digest.
ARG VERSION=0.8.3-slim@sha256:21532405da3e974a878335bd5734008f93c6066ac99dddda47e474cdc67a6351

FROM ghcr.io/vectorize-io/hindsight:${VERSION}

# --- local-ml deps -----------------------------------------------------------
# Mirror hindsight-api-slim's `local-ml` extra. Install into the image venv
# explicitly: the slim venv was created by `uv sync` and ships no pip, so a
# bare install would land in user site-packages invisible to the runtime python.
#
# torch is installed FIRST from the CPU wheel index. Vector inference is
# in-process on CPU, so the ~2GB of CUDA libraries pulled by PyPI's default
# torch build are pure dead weight. Installing torch alone from the cpu index
# guarantees the CPU build satisfies sentence-transformers' torch requirement
# in the next step (no CUDA build sneaks back in).
RUN uv pip install --python /app/api/.venv/bin/python --no-cache \
        --index-url https://download.pytorch.org/whl/cpu \
        'torch>=2.6.0'
RUN uv pip install --python /app/api/.venv/bin/python --no-cache \
        'sentence-transformers>=3.3.0' \
        'transformers>=4.53.0'

# --- model bake-in -----------------------------------------------------------
ARG EMBEDDER=BAAI/bge-m3
ARG RERANKER=BAAI/bge-reranker-v2-m3

# Cache weights under a fixed, world-readable path. The runtime runs non-root;
# without an explicit HF_HOME the weights would download into the build user's
# home and the runtime UID could not read them -> silent re-download at startup
# (and a hard failure if HF is unreachable through a cross-border proxy).
ENV HF_HOME=/app/models \
    HF_HUB_DOWNLOAD_TIMEOUT=600

RUN python -c "\
from sentence_transformers import SentenceTransformer, CrossEncoder; \
SentenceTransformer('${EMBEDDER}'); \
CrossEncoder('${RERANKER}')" \
    && chmod -R a+rX /app/models

# --- runtime model selection -------------------------------------------------
# Point Hindsight at the baked-in local models. These are defaults only; any
# value set in the pod spec / HelmRelease env overrides them.
ENV HINDSIGHT_API_EMBEDDINGS_PROVIDER=local \
    HINDSIGHT_API_EMBEDDINGS_LOCAL_MODEL=${EMBEDDER} \
    HINDSIGHT_API_RERANKER_PROVIDER=local \
    HINDSIGHT_API_RERANKER_LOCAL_MODEL=${RERANKER}

# --- CPU thread caps ---------------------------------------------------------
# torch/OpenMP otherwise grab every host core, ignoring the cgroup CPU quota and
# causing throttling + noisy-neighbour on shared nodes. Keep these aligned with
# resources.limits.cpu (4 here).
ENV OMP_NUM_THREADS=4 \
    MKL_NUM_THREADS=4

RUN rm -rf /var/cache/* /tmp/* /root/.cache 2>/dev/null || true

EXPOSE 8888 9999
