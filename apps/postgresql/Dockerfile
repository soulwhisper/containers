# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# postgresql: official postgres + the extensions the cnpg stack needs —
# the cloudnative-pg operator baseline (pgaudit, pg_failover_slots, pgvector,
# PostGIS, uuid-ossp, jq) plus this cluster's search extensions (vchord,
# vchord_bm25, pg_tokenizer, pgroonga) — all installed via the pig CLI from
# the pigsty extension catalog.
#
# Versioning model: extension versions are NOT pinned — every rebuild
# (triggered by an upstream postgres update via the BASE digest) installs
# the latest pig CLI and latest PG-major-compatible extensions. Released
# with the postgres major.minor.patch tag (e.g. :18.6) so cnpg clusters pin
# a reproducible version.
#
# WAL archiving/backups are handled by the plugin-barman-cloud deployment
# (cnpg's external barman plugin) — no barman binaries needed in this image.
#
# Why not upstream images:
#   ghcr.io/tensorchord/cloudnative-vectorchord:18 ships only vchord.so —
#   no vchord_bm25 (preload died in cnpg initdb), no pgroonga.
ARG BASE=18.6-bookworm@sha256:1c59e2c3c818eaa0f0628f695b36e7c9e362d6b219b36a54a32df645cbd7e1af
ARG PG_MAJOR


FROM postgres:${BASE}

# re-declared after FROM (pre-FROM args do not persist into build stages)
ARG PG_MAJOR

# install pig (latest pigsty CLI) -> all extensions via pig (latest compatible)
RUN <<EOF
  set -e
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates jq
  curl -fsSL https://repo.pigsty.io/pig | bash
  rm -f /etc/apt/sources.list.d/pgdg.sources /etc/apt/sources.list.d/pgdg.list  # pig manages pgdg sources itself
  pig repo add pgsql -u
  pig ext install -y vchord vchord_bm25 pg_tokenizer pgroonga uuid-ossp pgaudit pg_failover_slots vector postgis -v ${PG_MAJOR}
  apt-get purge -y --auto-remove curl
  rm -rf /var/lib/apt/lists/*
EOF

# cnpg runs instance containers as uid 26 (plain postgres image uses 999);
# initdb fails with 'user does not exist' without this
RUN usermod -u 26 postgres
USER 26

CMD ["postgres", "-c", "shared_preload_libraries=vchord,vchord_bm25,vector,pg_tokenizer,pgroonga", "-c", "search_path=\"$user\", public, bm25_catalog, tokenizer_catalog"]
