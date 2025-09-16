## BIRD in alpine

- The BIRD project is a fully functional dynamic IP routing daemon primarily targeted on (but not limited to) Linux.
- This image run BIRD on alpine, which is smaller and faster.
- Alternative: frr container, official repo, also alpine, [link](https://quay.io/repository/frrouting/frr?tab=tags).

### Example Usage

- cilium feature BFD is behind paid gate, [ref](https://github.com/cilium/cilium/issues/22394);
- using BIRD to provide BFD support;

```yaml
# manifest
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: bird
  namespace: kube-system
  labels:
    k8s-app: bird
    app.kubernetes.io/name: bird
spec:
  selector:
    matchLabels:
      k8s-app: bird
  template:
    metadata:
      labels:
        k8s-app: bird
        app.kubernetes.io/name: bird
    spec:
      containers:
        - name: alpine-bird
          image: "ghcr.io/soulwhisper/alpine-bird:edge"
          volumeMounts:
            - name: config
              mountPath: /etc/bird.conf
              subPath: bird.conf
              readOnly: true
      restartPolicy: Always
      hostNetwork: true
      volumes:
        - name: config
          configMap:
            name: bird-config
```

```conf
# bird.conf
## required to autoconfigure router-id
protocol device {}

## bare-bones BFD configuration
protocol bfd {
    interface "bond0" {
        interval 300 ms;
        multiplier 3;
    };
    neighbor 10.10.0.2; # bgp router
}
```
