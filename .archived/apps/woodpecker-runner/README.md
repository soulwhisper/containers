# woodpecker-runner

Step runner image for the home-ops `.woodpecker/` pipelines. Bakes every tool
the pipelines would otherwise download per run (crane, kubeconform,
dnscontrol, shellcheck, vendir, talosctl, yq, flate via mise; git, curl, jq,
python3+pyyaml via apt) so steps start without install warm-up.

## Base choice: debian, not alpine

glibc base on purpose — musl gaps for CI tooling:

- python wheels without musllinux builds fall back to slow source builds
- some release binaries ship glibc-only builds
- docker CLI plugin lookup paths differ under busybox/musl

Steps that need no tool installs keep their official images instead
(`docker:28-cli` for compose validation, `ghcr.io/astral-sh/uv` for docs,
`ghcr.io/cloudnative-pg/postgresql` for barman checks).

## Versions

Tool versions live in `app/.mise.toml` (renovate-tracked). `talosctl` must
match the cluster's Talos major.minor — keep aligned with home-ops
`.mise.toml`.
