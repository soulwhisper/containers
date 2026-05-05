# syntax=docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769

FROM ghcr.io/catthehacker/ubuntu:act-24.04@sha256:eb153df03ca81e75f73d40a2b14b5dbd2605e7813556efab47b23dfbfceeaa9d

# Install packages from other images
COPY --from=ghcr.io/fluxcd/flux-cli:v2.8.6@sha256:8dac34106b9d684ce07b13c9762b5132813c3d477cfd724e9e13cc5944ad771e /usr/local/bin/flux /usr/local/bin/flux
COPY --from=docker.io/alpine/helm:4.1.1@sha256:b5c85b997d83e89d9e8ff9215a14b03864274143981af45eb3fe729cdf782c73 /usr/bin/helm /usr/local/bin/helm
COPY --from=ghcr.io/sigstore/cosign/cosign:v3.0.6@sha256:de9c65609e6bde17e6b48de485ee788407c9502fa08b8f4459f595b21f56cd00 /ko-app/cosign /usr/local/bin/cosign

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
