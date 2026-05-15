# syntax=docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769

# Hermes Suite — all-in-one image combining hermes-agent, hermes-webui, and
# the built-in hermes dashboard, supervised by a single supervisord.
#
# Service ports:
#   8642 — hermes-gateway   (agent CLI / telegram / cron / tools)
#   9119 — hermes-dashboard (analytics, built into hermes-agent)
#   8787 — hermes-webui     (browser chat UI)

ARG AGENT_VERSION=v2026.5.7@sha256:9462082533e603fc6ba4259c51d7722dcd8cbb0f348488ceedc35d06fa8f9324

FROM nousresearch/hermes-agent:${AGENT_VERSION}

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

# Tinker-Atropos reasoning engine and browser tool dependencies.
# uv is already present in the base image.
RUN uv pip install -e /opt/hermes/tinker-atropos \
    && cd /opt/hermes \
    && npm install --prefer-offline --no-audit \
    && npx playwright install --with-deps chromium \
    && rm -rf /opt/hermes/scripts/whatsapp-bridge /var/lib/apt/lists/*

# supervisord is not packaged for Debian Trixie; install it into its own venv.
RUN uv venv /opt/supervisor \
    && uv pip install --python /opt/supervisor/bin/python3 supervisor \
    && ln -sf /opt/supervisor/bin/supervisord /usr/local/bin/supervisord \
    && ln -sf /opt/supervisor/bin/supervisorctl /usr/local/bin/supervisorctl

# hermes-webui — pinned to a release tag for reproducibility.
# Shares the agent's installed package so it can import agent modules.
RUN cd /opt \
    && git clone --depth 1 --branch "${WEBUI_VERSION}" \
        https://github.com/nesquena/hermes-webui.git hermes-webui \
    && uv pip install --python /opt/hermes/.venv/bin/python3 --no-cache-dir \
        -r /opt/hermes-webui/requirements.txt \
    && rm -rf /opt/hermes-webui/.git \
    && echo "__version__ = '${WEBUI_VERSION}'" > /opt/hermes-webui/api/_version.py

ENV PATH="/opt/hermes/.venv/bin:$PATH" \
    VIRTUAL_ENV=/opt/hermes/.venv \
    HERMES_HOME=/opt/data \
    HERMES_DATA_DIR=/opt/data \
    HOME=/opt/data \
    PLAYWRIGHT_BROWSERS_PATH=/opt/hermes/.playwright \
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
    && chown -R hermes:hermes /opt/data /workspace /opt/hermes-webui

# Cleanup
RUN rm -rf /var/cache/* /var/log/apt/* /tmp/* /root/.cache

EXPOSE 8642 8787 9119

WORKDIR /opt/hermes

USER hermes

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["supervisord", "-c", "/app/supervisord.conf", "-n"]
