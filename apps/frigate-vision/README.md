## frigate-vision

Python runtime image for the home-ops `smarthome-apps/frigate-vision` MQTT→hermes forwarder, with `paho-mqtt` and `httpx` baked in at build time. VLM description runs in frigate's native GenAI; this image only carries the forwarder's deps.

Why this exists: the cluster pod previously ran `pip install paho-mqtt httpx` at every boot, which requires Internet egress that the pod's CiliumNetworkPolicy deliberately does not grant. Baked-in deps remove the runtime pypi.org dependency entirely.

Dependency versions are pinned as `ARG`s in [`docker-bake.hcl`](./docker-bake.hcl) and tracked by renovate (pypi datasource).

### Discovery

The default `CMD` prints the bundled dependency versions:

```sh
docker run --rm ghcr.io/soulwhisper/frigate-vision:latest
```

### Usage

```yaml
containers:
  app:
    image:
      repository: ghcr.io/soulwhisper/frigate-vision
      tag: "<date-tag>"
    command: ["python", "/app/bridge.py"]
```
