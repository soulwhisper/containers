# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

# Hermes Suite — all-in-one image combining hermes-agent, hermes-webui, and
# the built-in hermes dashboard, supervised by a single supervisord.
#
# Service ports:
#   8642 — hermes-gateway   (agent CLI / telegram / cron / tools)
#   9119 — hermes-dashboard (analytics, built into hermes-agent)
#   8787 — hermes-webui     (browser chat UI)

ARG AGENT_VERSION=v2026.6.5@sha256:9ad3b04ec916ea2c2da22358fd43b024c788d74073210695af88bfc2e63869b4

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
    && chown -R hermes:hermes /opt/data /workspace /opt/hermes-webui \
    && echo 'export PATH="$HOME/.local/bin:$PATH"' >> $HOME/.bashrc

# Cleanup
RUN rm -rf /var/cache/* /var/log/apt/* /tmp/* /root/.cache

EXPOSE 8642 8787 9119

WORKDIR /opt/hermes

USER hermes

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["supervisord", "-c", "/app/supervisord.conf", "-n"]
