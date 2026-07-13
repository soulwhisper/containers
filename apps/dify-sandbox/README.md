## dify-sandbox

Extended [Dify sandbox](https://github.com/langgenius/dify-sandbox) image with **mise**, **node**, and **officecli** pre-installed.

Based on `langgenius/dify-sandbox:main` — the upstream sandbox server (`/main`) runs by default.

### Additions over upstream

| Layer | Detail |
|-------|--------|
| mise | Installed via apt; shims on system PATH (`/usr/local/mise/shims`) |
| node | mise built-in backend (takes precedence over upstream's extracted node) |
| officecli | `npm:@officecli/officecli` via mise npm backend |

A non-root `sandbox` user is available for interactive exec sessions.

### Example usage

```yaml
spec:
  containers:
    - name: sandbox
      image: ghcr.io/soulwhisper/dify-sandbox:latest
```

```sh
# exec as sandbox user to use officecli
kubectl exec -it deploy/dify-sandbox -- su - sandbox
kubectl exec -it deploy/dify-sandbox -- su - sandbox -c 'officecli --help'
```
