# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

ARG BASE=18.6-bookworm@sha256:1c59e2c3c818eaa0f0628f695b36e7c9e362d6b219b36a54a32df645cbd7e1af

FROM postgres:${BASE}
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
EOF

# cnpg runs instance containers as uid 26
# ref:https://github.com/cloudnative-pg/postgres-containers/blob/main/Dockerfile
RUN usermod -u 26 postgres
USER 26

CMD ["postgres", "-c", "shared_preload_libraries=vchord,vchord_bm25,vector,pg_tokenizer,pgroonga", "-c", "search_path=\"$user\", public, bm25_catalog, tokenizer_catalog"]
