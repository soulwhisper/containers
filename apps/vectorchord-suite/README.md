# vectorchord-suite

PostgreSQL 18 (official `postgres` image) with the VectorChord extension
suite and pgroonga, installed via the [pig](https://pig.pgsty.com) package
manager from the pigsty extension catalog.

## Extensions (PG 18)

| Extension      | Source (pigsty catalog)                | Used by                                |
| -------------- | -------------------------------------- | -------------------------------------- |
| `vchord`       | `postgresql-18-vchord`                 | immich, hindsight (vector search)      |
| `pgroonga`     | `postgresql-18-pgroonga`               | hindsight (keyword search)             |
| `vchord_bm25`  | `postgresql-18-vchord-bm25`            | optional BM25 hybrid search            |
| `pg_tokenizer` | `postgresql-18-pg-tokenizer`           | tokenizer required by `vchord_bm25`    |
| `vector`       | `postgresql-18-pgvector` (pgvector)    | general purpose                        |

## Why not upstream images

- `ghcr.io/tensorchord/cloudnative-vectorchord:18` ships only `vchord.so` —
  no `vchord_bm25` (cnpg initdb died on the preload) and no `pgroonga`.
- `ghcr.io/tensorchord/cloudnative-vectorchord-bm25` is not publicly pullable
  (ghcr 403).

## Versioning

Extension versions are **not** pinned: `pig ext install -v 18` resolves the
latest compatible builds per PG major — the same float model upstream images
use. Only two things are pinned:

- the postgres base image digest (`18-bookworm@sha256:…`, renovate-tracked)
- the pig CLI version (`PIG_VERSION` build arg, renovate-tracked in
  `docker-bake.hcl`)

`mise` is intentionally not used to install pig: `github:pgsty/pig` is not in
the mise registry and the github backend hits unauthenticated API rate limits
in CI builds. The pig CLI is installed via the official pigsty installer
(`curl -fsSL https://repo.pigsty.io/pig | bash -s ${PIG_VERSION}`).

## Preload

The cnpg Cluster spec owns `shared_preload_libraries`; the standalone CMD
(also kept from upstream) uses:

```
vchord,vchord_bm25,vector,pg_tokenizer,pgroonga
```
