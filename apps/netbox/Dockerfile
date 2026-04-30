# syntax=docker/dockerfile:1

ARG VERSION=latest
FROM netboxcommunity/netbox:${VERSION}

# Install NetBox plugins into the existing venv via uv (already shipped in the base image)
RUN . /opt/netbox/venv/bin/activate \
    && uv pip install \
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
RUN rm -rf /var/cache/* /var/log/* /tmp/* /root/.cache
