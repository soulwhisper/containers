## manyfold — non-root solo image

Built from [`manyfold3d/manyfold`](https://github.com/manyfold3d/manyfold)
at the version selected in `docker-bake.hcl`, including its bundled Redis. A
custom entrypoint and `supervisord` run Redis and Manyfold directly; the
upstream s6-overlay entrypoint is unused.

The image supports caller-selected non-root UID/GID execution and a read-only
root filesystem. Runtime identity and environment belong to the deployment;
no root entrypoint or runtime capabilities are required.

### Required runtime identity

Set `PUID` and `PGID` explicitly to the numeric UID/GID used to run the
container, such as the pod's `runAsUser` and `runAsGroup`. Startup uses these
values for SQLite database and plugin-directory ownership. The inherited
upstream default is `0` for both; omitting them from a non-root deployment can
fail or assign root ownership to application files.

### Required writable mounts

- `/config` — SQLite database, settings, and plugins
- `/tmp` — private temporary and supervisor state
- `/usr/src/app/tmp` — application temporary files
- `/usr/src/app/log` — application logs
- `/libraries/main` — model library storage

`SECRET_KEY_BASE` must be set. The image is published as
`ghcr.io/soulwhisper/manyfold:<upstream-version>`.
