# cloudnative-custom

PostgreSQL (official `postgres` image) with the CloudNativePG operator's
baseline extensions plus this cluster's search extensions, installed via the
[pig](https://pig.pgsty.com) package manager from the pigsty extension catalog.

## Extensions

| Extension            | Used by                                        |
| -------------------- | ---------------------------------------------- |
| `vchord`             | immich, hindsight (vector search)              |
| `pgroonga`           | hindsight (keyword search)                     |
| `vchord_bm25`        | optional BM25 hybrid search                    |
| `vector` (pgvector)  | **required by vchord** (package + SQL dependency) |
| `pgaudit`            | cnpg baseline                                  |
| `pg_failover_slots`  | cnpg baseline                                  |
The baseline set mirrors the official `cloudnative-pg/postgresql` bundle

## Why not upstream images

- `ghcr.io/tensorchord/cloudnative-vectorchord:18` ships only `vchord.so` —
  no `vchord_bm25` (cnpg initdb died on the preload) and no `pgroonga`.

## Versioning

- Extension versions are **not** pinned: every rebuild installs the latest
  pig CLI and the latest PG-major-compatible extension builds — the same
  float model upstream images use.
- Builds are **triggered by upstream postgres updates** only: the `BASE`
  digest is renovate-tracked, so a new upstream postgres release produces a
  rebuild.
- The image is released with the **postgres major tag** (e.g. `:18`).

`mise` is intentionally not used to install pig: `github:pgsty/pig` is not in
the mise registry and the github backend hits unauthenticated API rate limits
in CI builds. The pig CLI is installed via the official pigsty installer
(`curl -fsSL https://repo.pigsty.io/pig | bash`).

## Preload

The cnpg Cluster spec owns `shared_preload_libraries`; the standalone CMD
(also kept from upstream) uses:

```
vchord,vchord_bm25,vector,pg_tokenizer,pgroonga
```
