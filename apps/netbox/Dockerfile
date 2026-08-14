# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

ARG VERSION=v4.6.8@sha256:1abb7a50280cde87aca7e1a7391cde620f18c617333b24692adb1c077e58ecda

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
