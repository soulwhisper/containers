# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# postgresql: official postgres + the extensions the cnpg stack
# needs — the cloudnative-pg operator baseline (pgaudit, pg_failover_slots,
# pgvector, PostGIS) plus this cluster's search extensions (vchord,
# vchord_bm25, pg_tokenizer, pgroonga) — all installed via the pig CLI
# from the pigsty extension catalog.
# Versioning model: extension versions are NOT pinned — every rebuild
# (triggered by an upstream postgres update via the BASE digest) installs
# the latest pig CLI and latest PG-major-compatible extensions. Only the
# postgres base digest is renovate-tracked; the image is released with the
# postgres major tag (e.g. :18).
#
# Why not upstream images:
#   ghcr.io/tensorchord/cloudnative-vectorchord:18 ships only vchord.so —
#   no vchord_bm25 (preload died in cnpg initdb), no pgroonga.
#   ghcr.io/tensorchord/cloudnative-vectorchord-bm25 is not publicly pullable.
ARG BASE=18.6-bookworm@sha256:1c59e2c3c818eaa0f0628f695b36e7c9e362d6b219b36a54a32df645cbd7e1af

FROM postgres:${BASE}

# install pig (latest pigsty CLI) -> all extensions via pig (latest compatible)
RUN <<EOF
  set -e
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates
  curl -fsSL https://repo.pigsty.io/pig | bash
  rm -f /etc/apt/sources.list.d/pgdg.list  # pig manages pgdg sources itself
  pig repo add pgsql -u
  pig ext install -y vchord vchord_bm25 pg_tokenizer pgroonga vector pgaudit pg_failover_slots postgis -v ${PG_MAJOR}
  # NOTE: pgvector is NOT optional — postgresql-18-vchord Depends on it and
  # CREATE EXTENSION vchord fails without it (verified). Do not drop it.
  apt-get purge -y --auto-remove curl
  rm -rf /var/lib/apt/lists/*
EOF
CMD ["postgres", "-c", "shared_preload_libraries=vchord,vchord_bm25,vector,pg_tokenizer,pgroonga", "-c", "search_path=\"$user\", public, bm25_catalog, tokenizer_catalog"]
