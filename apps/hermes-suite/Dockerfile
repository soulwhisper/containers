# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# Hermes Suite — all-in-one image combining hermes-agent, hermes-webui, and
# the built-in hermes dashboard, supervised by a single supervisord.
#
# Service ports:
#   8642 — hermes-gateway   (agent CLI / telegram / cron / tools)
#   9119 — hermes-dashboard (analytics, built into hermes-agent)
#   8787 — hermes-webui     (browser chat UI)

# Base: upstream `stable` channel (currently v0.21.7), digest-pinned.
# Docker Hub publishes no plain vX.Y.Z tag — only immutable rc.N-vX.Y.Z
# candidates and the floating stable/latest aliases — so renovate floats
# this digest forward on each upstream release while the suite's version
# label/tag is tracked from GitHub releases in docker-bake.hcl.
FROM nousresearch/hermes-agent:stable@sha256:444b42594e7f467425872fc35e077a227d377f0fe0726bc818706d62175fad8b

ARG WEBUI_VERSION

USER root

# Runtime utilities + gosu for clean privilege drop in the entrypoint.
RUN apt-get -qq update \
    && apt-get -qq install --no-install-recommends \
        git \
        curl \
        nano \
        net-tools \
        iputils-ping \
        iproute2 \
        openssh-client \
        procps \
        gosu \
    && rm -rf /var/lib/apt/lists/*

# supervisord is not packaged for Debian Trixie; install it into its own venv.
# Bind to the Debian system python (/usr/bin/python3): /usr/local/bin/python3
# is upstream's pm-managed interpreter (under /opt/hermes/tools); a venv bound
# to it would break if upstream relocates the store. Upstream no longer ships
# uv on PATH (pm-managed since v0.21.6), so plain venv + pip.
RUN /usr/bin/python3 -m venv /opt/supervisor \
    && /opt/supervisor/bin/pip install --no-cache-dir supervisor \
    && ln -sf /opt/supervisor/bin/supervisord /usr/local/bin/supervisord \
    && ln -sf /opt/supervisor/bin/supervisorctl /usr/local/bin/supervisorctl

# hermes-webui — pinned to a release tag for reproducibility.
# Shares the agent's installed package so it can import agent modules.
RUN cd /opt \
    && git clone --depth 1 --branch "${WEBUI_VERSION}" \
        https://github.com/nesquena/hermes-webui.git hermes-webui \
    && /opt/hermes/.venv/bin/python3 -m ensurepip \
    && /opt/hermes/.venv/bin/python3 -m pip install --no-cache-dir \
        -r /opt/hermes-webui/requirements.txt \
    && rm -rf /opt/hermes-webui/.git \
    && echo "__version__ = '${WEBUI_VERSION}'" > /opt/hermes-webui/api/_version.py

# Platform plugins — deps pinned to upstream extras (plugins/platforms/*/).
# Baked in because runtime lazy-install fails under the cluster egress policy.
# feishu = ["lark-oapi==1.6.8", "qrcode==7.4.2"] (pyproject.toml)
# The project venv has pip seeded by the webui install above.
RUN /opt/hermes/.venv/bin/python3 -m pip install --no-cache-dir \
        "lark-oapi==1.6.8" \
        "qrcode==7.4.2"

# PLAYWRIGHT_BROWSERS_PATH intentionally not overridden: upstream ships
# pm-managed Chromium under /opt/hermes/tools and sets the path itself.
ENV PATH="/opt/hermes/.venv/bin:$PATH" \
    VIRTUAL_ENV=/opt/hermes/.venv \
    HERMES_HOME=/opt/data \
    HERMES_DATA_DIR=/opt/data \
    HOME=/opt/data \
    HERMES_WEB_DIST=/opt/hermes/hermes_cli/web_dist \
    HERMES_WEBUI_HOST=0.0.0.0 \
    HERMES_WEBUI_PORT=8787 \
    HERMES_WEBUI_STATE_DIR=/opt/data/webui \
    HERMES_WEBUI_DEFAULT_WORKSPACE=/workspace \
    HERMES_WEBUI_AGENT_DIR=/opt/hermes \
    HERMES_SERVICES=gateway,webui

# Entrypoint
COPY app /app
RUN chmod +x /app/entrypoint.sh \
    && mkdir -p /workspace /opt/data \
    && chown -R hermes:hermes /opt/data /workspace /opt/hermes-webui \
    && echo 'export PATH="$HOME/.local/bin:$PATH"' >> $HOME/.bashrc

# Cleanup
RUN rm -rf /var/cache/* /var/log/apt/* /tmp/* /root/.cache

EXPOSE 8642 8787 9119

WORKDIR /opt/hermes

USER hermes

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["supervisord", "-c", "/app/supervisord.conf", "-n"]
