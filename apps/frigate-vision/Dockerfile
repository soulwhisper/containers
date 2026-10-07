# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# frigate-vision bridge runtime: python with paho-mqtt + httpx baked in.
# Exists so the cluster pod never pip-installs at boot — the frigate-vision
# CiliumNetworkPolicy allows no Internet egress, and runtime installs are a
# pypi.org availability dependency either way.

FROM docker.io/library/python:3.14-alpine3.23@sha256:e7cff362a12454395f8dc99607c21a602ea880cdce03538877b90c2ba932b16f

ARG PAHO_MQTT_VERSION
ARG HTTPX_VERSION

RUN pip install --no-cache-dir --no-compile \
    "paho-mqtt==${PAHO_MQTT_VERSION}" \
    "httpx==${HTTPX_VERSION}"

# Discovery: print the bundled dependency table (no side effects).
CMD ["python", "-c", "import importlib.metadata as m; [print(f'{d}=={m.version(d)}') for d in ('paho-mqtt', 'httpx')]"]
