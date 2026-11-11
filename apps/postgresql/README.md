# postgresql

PostgreSQL (official `postgres` image) with the CloudNativePG operator's
baseline extensions, PostgreSQL contrib, and this cluster's search extensions.
The non-contrib extensions are installed via the
[pig](https://pig.pgsty.com) package manager from the pigsty extension catalog.

## Extensions

| Extension            | Source             | Used by                                      |
| -------------------- | ------------------ | -------------------------------------------- |
| `pgroonga`           | pig                | hindsight (keyword search)                   |
| `vector` (pgvector)  | pig                | hindsight, open-webui; immich pgvector restore |
| `uuid-ossp`          | pig                | immich (UUID generation)                     |
| `unaccent`           | postgres contrib   | immich (search)                              |
| `cube`               | postgres contrib   | immich (`earthdistance` dependency)          |
| `earthdistance`      | postgres contrib   | immich (geo search)                          |
| `pg_trgm`            | postgres contrib   | immich (trigram search)                      |
| `pgaudit`            | pig                | cnpg baseline                                |
| `pg_failover_slots`  | pig                | cnpg baseline                                |
| `postgis`            | pig                | available; no active consumer                |

`CREATE EXTENSION` and CNPG `Database` reconciliation only register extension
files that are already present in the image; they do not install a missing
package. `CASCADE` resolves extension dependencies, not OS packages.

## Why not upstream images

- The official CloudNativePG image includes the baseline and pgvector, but not
  `pgroonga`.

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
vector,pgroonga
```
