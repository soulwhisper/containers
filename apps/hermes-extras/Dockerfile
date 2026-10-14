# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# hermes-agent plugin deps (feishu platform: lark-oapi; requests-toolbelt is
# its only dependency missing from the upstream venv). --no-deps: everything
# else in the closure (requests, httpx, websockets, pycryptodome, qrcode,
# pypng) already exists in the upstream image venv — shipping copies here
# would shadow them via PYTHONPATH. install.sh mirrors /pylib onto the PVC.
FROM docker.io/library/python:3.14-alpine3.23@sha256:e7cff362a12454395f8dc99607c21a602ea880cdce03538877b90c2ba932b16f

ARG LARK_OAPI_VERSION
ARG REQUESTS_TOOLBELT_VERSION
RUN pip install --no-cache-dir --target=/pylib --no-deps \
      "lark-oapi==${LARK_OAPI_VERSION}" \
      "requests-toolbelt==${REQUESTS_TOOLBELT_VERSION}"

COPY app /app
RUN chown -R 10000:10000 /pylib /app

USER 10000:10000

# No ENTRYPOINT: init-containers and tests override `command` fully.
# Default CMD lists the bundle — useful for `docker run` discovery.
CMD ["ls", "-l", "/pylib"]
