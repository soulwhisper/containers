## Netbox Custom

- Tested with [nix-config](https://github.com/soulwhisper/nix-config/blob/main/hosts/_modules/nixos/services/netbox/default.nix);

### Plugins

- [netbox-attachments](https://github.com/Kani999/netbox-attachments)
- [netbox-bgp](https://github.com/netbox-community/netbox-bgp)
- [netbox-dns](https://github.com/peteeckel/netbox-plugin-dns)
- [netbox-floorplan-plugin](https://github.com/netbox-community/netbox-floorplan-plugin)
- [netbox-interface-synchronization](https://github.com/NetTech2001/netbox-interface-synchronization)
- [netbox-napalm-plugin](https://github.com/netbox-community/netbox-napalm-plugin)
- [netbox-plugin-prometheus-sd](https://github.com/FlxPeters/netbox-plugin-prometheus-sd)
- [netbox-qrcode](https://github.com/netbox-community/netbox-qrcode)
- [netbox-reorder-rack](https://github.com/netbox-community/netbox-reorder-rack)
- [netbox-topology-views](https://github.com/netbox-community/netbox-topology-views)

### Usage

- Deployed with official helm chart, values below

```yaml
image:
  repository: YOUR_REPOSITORY
  tag: IMAGE_TAG
plugins:
  - netbox_attachments
  - netbox_bgp
  - netbox_floorplan_plugin
  - netbox_interface_synchronization
  - netbox_plugin_dns
  - netbox_plugin_prometheus_sd
  - netbox_qrcode
  - netbox_reorder_rack
  - netbox_topology_views
#   - netbox_napalm_plugin
# pluginsConfig:
#   netbox_napalm_plugin:
#     NAPALM_USERNAME:
#       valueFrom:
#         secretKeyRef:
#           name: netbox-secrets
#           key: NAPALM_USERNAME
#     NAPALM_PASSWORD:
#       valueFrom:
#         secretKeyRef:
#           name: netbox-secrets
#           key: NAPALM_PASSWORD
```
