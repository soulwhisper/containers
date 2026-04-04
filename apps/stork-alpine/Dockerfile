FROM alpine:latest

COPY app /app

RUN chmod +x /app/entrypoint.sh && \
    mkdir -p /var/log/supervisor

RUN apk add --no-cache bash supervisor curl && \
    curl -1sLf 'https://dl.cloudsmith.io/public/isc/stork/setup.alpine.sh' | bash && \
    apk add --no-cache isc-stork-server isc-stork-agent

EXPOSE 8080 9547 9119

ENTRYPOINT ["/app/entrypoint.sh"]
