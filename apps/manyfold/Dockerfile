# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

ARG VERSION
FROM ghcr.io/manyfold3d/manyfold-solo:${VERSION}

RUN apk add --no-cache supervisor

COPY entrypoint.sh /usr/local/bin/manyfold-entrypoint.sh
COPY app/start.sh /usr/local/bin/manyfold-start.sh
COPY app/supervisord.conf /etc/supervisord.conf
RUN chmod 0555 /usr/local/bin/manyfold-entrypoint.sh /usr/local/bin/manyfold-start.sh

ENTRYPOINT ["/usr/local/bin/manyfold-entrypoint.sh"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf", "-n"]
