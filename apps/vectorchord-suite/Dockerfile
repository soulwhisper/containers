# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# VectorChord suite + pgroonga on official postgres.
#
# Toolchain follows the repo convention: apt-installed mise fetches the pig
# CLI (github:pgsty/pig release), and pig installs all PG extensions for
# PG 18 from the pigsty catalog — vchord, vchord_bm25, pg_tokenizer,
# pgroonga, vector (pgvector). Version of the pig CLI itself is pinned by
# PIG_VERSION (renovate-tracked in docker-bake.hcl); extension versions are
# resolved by pig per PG major.
#
# Why not upstream images:
#   ghcr.io/tensorchord/cloudnative-vectorchord:18 ships only vchord.so —
#   no vchord_bm25 (preload died in cnpg initdb), no pgroonga.
#   ghcr.io/tensorchord/cloudnative-vectorchord-bm25 is not publicly pullable.

ARG BASE=18-bookworm@sha256:1c59e2c3c818eaa0f0628f695b36e7c9e362d6b219b36a54a32df645cbd7e1af

FROM postgres:${BASE}

ARG PIG_VERSION

# install pig (pigsty CLI, version-pinned) -> all extensions via pig.
# NOTE: mise is not used for pig here — github:pgsty/pig is not in the mise
# registry and the github backend hits unauthenticated API rate limits in CI.
RUN <<EOF
  set -e
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates
  curl -fsSL https://repo.pigsty.io/pig | bash -s ${PIG_VERSION}
  rm -f /etc/apt/sources.list.d/pgdg.list  # pig manages pgdg sources itself
  pig repo add pgsql -u
  pig ext install -y vchord vchord_bm25 pg_tokenizer pgroonga vector -v 18
  apt-get purge -y --auto-remove curl
  rm -rf /var/lib/apt/lists/*
EOF

CMD ["postgres", "-c", "shared_preload_libraries=vchord,vchord_bm25,vector,pg_tokenizer,pgroonga", "-c", "search_path=\"$user\", public, bm25_catalog, tokenizer_catalog"]
