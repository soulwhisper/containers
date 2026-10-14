# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# hermes-agent — upstream nousresearch/hermes-agent plus build-time extras
# (rtk, caveman skills, feishu plugin deps) and a supervisord that runs
# gateway + dashboard in one container.
#   8642 — gateway (agent API / feishu / cron / webhooks)
#   9119 — dashboard (built-in web UI; auth gate requires a provider)

# Base: upstream `stable` channel (currently v0.21.7), digest-pinned.
# Docker Hub publishes no plain vX.Y.Z tag — only immutable rc.N-vX.Y.Z
# candidates and the floating stable/latest aliases — so renovate floats this
# digest forward on each upstream release while the image tag is tracked from
# GitHub releases in docker-bake.hcl; both updates ship in one grouped PR
# (minimumGroupSize=2).
FROM nousresearch/hermes-agent:stable@sha256:444b42594e7f467425872fc35e077a227d377f0fe0726bc818706d62175fad8b

USER root

# Runtime utilities + gosu for clean privilege drop in the entrypoint
# (git/curl/iputils-ping/openssh-client/procps already ship upstream).
RUN apt-get -qq update \
    && apt-get -qq install --no-install-recommends \
        nano \
        net-tools \
        iproute2 \
        gosu \
    && rm -rf /var/lib/apt/lists/*

# supervisord is not packaged for Debian Trixie; install it into its own venv
# bound to the Debian system python (upstream's python is pm-managed under
# /opt/hermes/tools; uv is no longer on PATH upstream).
RUN /usr/bin/python3 -m venv /opt/supervisor \
    && /opt/supervisor/bin/pip install --no-cache-dir supervisor \
    && ln -sf /opt/supervisor/bin/supervisord /usr/local/bin/supervisord \
    && ln -sf /opt/supervisor/bin/supervisorctl /usr/local/bin/supervisorctl

# Feishu platform plugin deps — pinned by UPSTREAM itself (pyproject
# optional-dependencies.feishu) and installed into the project venv at build
# time. Runtime lazy-install is impossible (read-only rootfs).
RUN set -eu; \
    mkdir -p /bundled; \
    /opt/hermes/.venv/bin/python3 -m ensurepip; \
    reqs="$(/opt/hermes/.venv/bin/python3 -c 'import tomllib; print(" ".join(tomllib.load(open("/opt/hermes/pyproject.toml","rb"))["project"]["optional-dependencies"]["feishu"]))')"; \
    echo "feishu extra pins (upstream): $reqs"; \
    /opt/hermes/.venv/bin/python3 -m pip install --no-cache-dir $reqs; \
    echo "feishu=$reqs" > /bundled/VERSIONS

# Build-time extras — FLOATING: resolved to latest at every build and never
# trigger a build on their own (only hermes-agent is version-tracked here;
# a new base digest/version rebuild picks up whatever extras are current).
RUN set -eux; \
    mkdir -p /bundled/bin /bundled/skills /bundled/plugins/rtk-rewrite; \
    cd /tmp; \
    # rtk — latest release binary, sha256-verified against the release's
    # own checksums.txt; the hermes plugin adapter ships in the git repo at
    # the same tag (static files — the runtime init's config mutation is
    # GitOps-managed, see configmap plugins.enabled)
    rtk_tag="$(curl -fsSL -o /dev/null -w '%{url_effective}' https://github.com/rtk-ai/rtk/releases/latest)"; \
    rtk_tag="${rtk_tag##*/}"; \
    case "$(uname -m)" in \
      x86_64)  rtk_asset="rtk-oss-x86_64-unknown-linux-musl" ;; \
      aarch64) rtk_asset="rtk-oss-aarch64-unknown-linux-gnu" ;; \
      *) echo "unsupported arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL "https://github.com/rtk-ai/rtk/releases/download/${rtk_tag}/checksums.txt" -o checksums.txt; \
    curl -fsSL "https://github.com/rtk-ai/rtk/releases/download/${rtk_tag}/${rtk_asset}.tar.gz" -o "${rtk_asset}.tar.gz"; \
    grep "  ${rtk_asset}.tar.gz\$" checksums.txt | sha256sum -c -; \
    tar -xzf "${rtk_asset}.tar.gz" -C /bundled/bin rtk; \
    chmod 755 /bundled/bin/rtk; \
    for f in __init__.py plugin.yaml; do \
      curl -fsSL "https://raw.githubusercontent.com/rtk-ai/rtk/${rtk_tag}/hooks/hermes/rtk-rewrite/${f}" \
        -o "/bundled/plugins/rtk-rewrite/${f}"; \
    done; \
    echo "rtk=${rtk_tag}" >> /bundled/VERSIONS; \
    # caveman — latest release tag resolved via the releases/latest redirect
    cave_tag="$(curl -fsSL -o /dev/null -w '%{url_effective}' https://github.com/JuliusBrussee/caveman/releases/latest)"; \
    cave_tag="${cave_tag##*/}"; \
    curl -fsSL "https://github.com/JuliusBrussee/caveman/archive/refs/tags/${cave_tag}.tar.gz" -o caveman.tar.gz; \
    tar -xzf caveman.tar.gz; \
    cp -r "caveman-${cave_tag#v}/skills/." /bundled/skills/; \
    echo "caveman=${cave_tag}" >> /bundled/VERSIONS; \
    rm -rf /tmp/*

ENV PATH="/bundled/bin:$PATH" \
    RTK_TELEMETRY_DISABLED=true \
    HOME=/opt/data

# Entrypoint + supervisord
COPY app /app
RUN chmod +x /app/entrypoint.sh \
    && mkdir -p /workspace \
    && chown -R hermes:hermes /opt/data /workspace /bundled

# Cleanup
RUN rm -rf /var/cache/* /var/log/apt/* /tmp/* /root/.cache

EXPOSE 8642 9119

WORKDIR /opt/hermes

USER hermes

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["supervisord", "-c", "/app/supervisord.conf", "-n"]
