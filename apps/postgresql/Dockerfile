# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

ARG BASE=18.7-bookworm@sha256:ece55f7868da684e400f82f1030a050718adbdd65c23aa500c486f6e2950dc0f

FROM postgres:${BASE}
ARG PG_MAJOR

# install pig (latest pigsty CLI) -> all extensions via pig (latest compatible)
RUN <<EOF
  set -e
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates jq
  curl -fsSL https://repo.pigsty.io/pig | bash
  rm -f /etc/apt/sources.list.d/pgdg.sources /etc/apt/sources.list.d/pgdg.list  # pig manages pgdg sources itself
  # silence pig's untrusted-staging-update noise: pgdg key ships with the base
  # image; pigsty's own key is unpublished upstream (hence pig's trusted=yes
  # default) so its NO_PUBKEY warning during 'repo add' is cosmetic
  cp /usr/local/share/keyrings/postgres.gpg.asc /etc/apt/trusted.gpg.d/postgresql.asc
  pig repo add pgsql -u
  pig ext install -y pgroonga pgaudit pg_failover_slots vector -v ${PG_MAJOR}
EOF

# cnpg runs instance containers as uid 26
# ref:https://github.com/cloudnative-pg/postgres-containers/blob/main/Dockerfile
RUN usermod -u 26 postgres
USER 26

CMD ["postgres", "-c", "shared_preload_libraries=vector,pgroonga", "-c", "search_path=\"$user\", public"]
