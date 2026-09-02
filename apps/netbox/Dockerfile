# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

ARG VERSION=v4.6.10@sha256:91b823a05cb51004f07acc2228ccc2993f38b0f0bf711b403b0fdf85e51277e8

FROM netboxcommunity/netbox:${VERSION}

# Install NetBox plugins into the existing venv via uv (already shipped in the base image)
ENV VIRTUAL_ENV=/opt/netbox/venv
RUN uv pip install \
    netbox-attachments \
    netbox-bgp \
    netbox-floorplan-plugin \
    netbox-interface-synchronization \
    netbox-plugin-dns \
    netbox-plugin-prometheus-sd \
    netbox-qrcode \
    netbox-reorder-rack \
    netbox-topology-views \
    netbox-napalm-plugin

# Cleanup
RUN rm -rf /var/cache/* /var/log/* /tmp/*
