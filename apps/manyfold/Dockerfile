# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

FROM ghcr.io/manyfold3d/manyfold-solo:0.154.1@sha256:f1b7784fb63aec500fca61081d3f970948fc652ad16ecdc5d57cd3c0c925127f

USER root
RUN apk add --no-cache supervisor

COPY entrypoint.sh /usr/local/bin/manyfold-entrypoint.sh
COPY start.sh /usr/local/bin/manyfold-start.sh
COPY supervisord.conf /etc/supervisord.conf
RUN chmod 0555 /usr/local/bin/manyfold-entrypoint.sh /usr/local/bin/manyfold-start.sh

ENV PUID=2000 \
    PGID=2000 \
    TMPDIR=/tmp/manyfold

USER 2000:2000
ENTRYPOINT ["/usr/local/bin/manyfold-entrypoint.sh"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf", "-n"]
