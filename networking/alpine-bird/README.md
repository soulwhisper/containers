## BIRD in alpine

- The BIRD project is a fully functional dynamic IP routing daemon primarily targeted on (but not limited to) Linux.
- This image run BIRD on alpine, which is smaller and faster.
- Alternative: frr container, official repo, also alpine, [link](https://quay.io/repository/frrouting/frr?tab=tags).

### Example Usage

- cilium feature BFD is behind paid gate, [ref](https://github.com/cilium/cilium/issues/22394);
- we are using BIRD to provide BFD support;

```shell
# manifest
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: bfd-bird
  namespace: kube-system
  labels:
    k8s-app: bfd-bird
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: cilium
      app.kubernetes.io/component: agent
  template:
    metadata:
      labels:
        app.kubernetes.io/name: cilium
        app.kubernetes.io/component: agent
    spec:
      containers:
      - image: "ghcr.io/container/alpine-bird:edge-2.16"
        name: bfd-bird
        volumeMounts:
        - mountPath: /etc/bird.conf
          name: config
          subPath: bird.conf
      hostNetwork: true
      volumes:
      - configMap:
          name: bird-config
        name: config

# bird.conf
## required to autoconfigure router-id
protocol device {}

## bare-bones BFD configuration
protocol bfd {
    interface "*" {
        interval 300 ms;
        multiplier 3;
    };
    neighbor 10.10.0.101;
    neighbor 10.10.0.102;
    neighbor 10.10.0.103;
}
```