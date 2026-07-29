# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

ARG VERSION=v4.6.6@sha256:974fb6ee8613fd75e236422f62e0ede82cd98cf5ac9f214a53dc889150cb39ec

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
