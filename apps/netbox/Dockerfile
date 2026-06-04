# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

ARG VERSION=v4.6.2@sha256:4f6284bff95fb53181c9c9f0561eb961d36044fa7484e25ff381c759cb54be6f

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
