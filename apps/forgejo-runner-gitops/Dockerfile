# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

FROM ghcr.io/catthehacker/ubuntu:act-24.04@sha256:f713414ec6c6f4c352c030005df8001c0c7d14154f3923e239d03daef2e9ae16

# Install packages from other images
COPY --from=ghcr.io/fluxcd/flux-cli:v2.9.4@sha256:5260c79fb1b744c78755d98bcb271971c93e4ea214623c3f9f96ff59536d0398 /usr/local/bin/flux /usr/local/bin/flux
COPY --from=docker.io/alpine/helm:4.1.1@sha256:b5c85b997d83e89d9e8ff9215a14b03864274143981af45eb3fe729cdf782c73 /usr/bin/helm /usr/local/bin/helm
COPY --from=ghcr.io/sigstore/cosign/cosign:v3.1.3@sha256:9e5c2f2edc34351160407ca3416c61855bdf9403c3c5936e0f0be7fc261611b8 /ko-app/cosign /usr/local/bin/cosign

# Install recent Python
RUN \
    add-apt-repository ppa:deadsnakes/ppa \
    && apt-get -qq update \
    && apt-get -qq install python3.14 python3.14-venv rsync

# Install kustomize
# renovate: datasource=github-releases depName=kubernetes-sigs/kustomize
ARG KUSTOMIZE_VERSION=v5.7.1
ARG TARGETARCH
RUN wget -qO- \
    "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_${TARGETARCH}.tar.gz" \
    | tar xz -C /usr/local/bin kustomize

# Install pipx
ENV PATH="/root/.local/bin/:$PATH"
RUN \
    mkdir -p ~/.local/bin \
    && python3.14 -m venv ~/.local/share/pipx-venv \
    && ~/.local/share/pipx-venv/bin/pip install pipx \
    && ln -s ~/.local/share/pipx-venv/bin/pipx ~/.local/bin/pipx

# Install flux-local
# renovate: datasource=github-releases depName=allenporter/flux-local
ARG FLUX_LOCAL_VERSION=8.1.0
RUN \
    pipx install flux-local==${FLUX_LOCAL_VERSION}

# Cleanup
RUN \
    apt-get clean && \
    rm -rf /var/cache/* /var/log/* /var/lib/apt/lists/* /tmp/*
