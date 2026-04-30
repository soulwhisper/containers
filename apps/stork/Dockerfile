# syntax=docker/dockerfile:1

FROM alpine:3.23

ARG VERSION

# Install runtime deps + ISC stork server & agent from upstream cloudsmith repo.
# Cloudsmith pkgver format is "<semver>.<timestamp>" (e.g. 2.4.0.260225XXXXXX),
# so we range-pin: floor at ${VERSION}, ceiling at the next minor line.
# This keeps builds reproducible while letting renovate bump VERSION on new releases.
RUN set -eux; \
    apk add --no-cache bash supervisor curl; \
    curl -1sLf 'https://dl.cloudsmith.io/public/isc/stork/setup.alpine.sh' | bash; \
    MAJOR="${VERSION%%.*}"; \
    MINOR_TMP="${VERSION#*.}"; \
    MINOR="${MINOR_TMP%%.*}"; \
    NEXT_MINOR="$((MINOR + 1))"; \
    apk add --no-cache \
        "isc-stork-server>=${VERSION}" "isc-stork-server<${MAJOR}.${NEXT_MINOR}.0" \
        "isc-stork-agent>=${VERSION}"  "isc-stork-agent<${MAJOR}.${NEXT_MINOR}.0"; \
    mkdir -p /var/log/supervisor

COPY app /app
RUN chmod +x /app/entrypoint.sh

# Cleanup
RUN rm -rf /var/cache/apk/* /tmp/* /root/.cache

EXPOSE 8080 9547 9119

ENTRYPOINT ["/app/entrypoint.sh"]
