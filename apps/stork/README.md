## ISC-Stork on alpine

Since the official ISC Stork project does not provide a native Docker image, this solution utilizes a lightweight Alpine Linux base to build a custom container. By integrating supervisord as the process manager and a dynamic entrypoint.sh script, this single image can be toggled between Server and Agent modes via environment variables. This approach ensures high availability with automatic process restarts while maintaining a minimal footprint and full compatibility with host-side services like Kea DHCP and BIND9.

### Example Usage

- server mode;

```yaml
---
services:
  stork-server:
    image: ghcr.io/soulwhisper/stork:latest
    container_name: stork-server
    restart: always
    ports:
      - "8080:8080"
    environment:
      - STORK_MODE=server
      - STORK_REST_PORT=8080
      - STORK_DATABASE_HOST=/run/postgresql
      - STORK_DATABASE_NAME=stork
      - STORK_DATABASE_USER_NAME=stork
      - STORK_SERVER_ENABLE_METRICS=1
    volumes:
      - /run/postgresql:/run/postgresql # host postgres socket
```

- agent mode;

```yaml
services:
  stork-agent:
    image: ghcr.io/soulwhisper/stork:latest
    container_name: stork-agent
    restart: always
    pid: "host"
    cap_add:
      - SYS_PTRACE
    ports:
      - "8081:8081/tcp"
      - "9119:9119/tcp" # bind9 metrics
      - "9547:9547/tcp" # kea metrics
    environment:
      - STORK_MODE=agent
      - STORK_AGENT_PORT=8081
      # docker use 'host.docker.internal', or 'http://stork-server:8080'.
      - STORK_AGENT_SERVER_URL=http://host.containers.internal:8080
    volumes:
      - /etc/bind:/etc/bind
      - /var/lib/kea:/var/lib/kea
      - /var/run/kea:/var/run/kea
```

- use env file;

```yaml
services:
  stork-server:
    image: ghcr.io/soulwhisper/stork:2.4.0
    container_name: stork-server
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - ./stork-server.env:/etc/stork/server.env:ro
      - /run/postgresql:/run/postgresql
```

```
# stork-server.env
STORK_REST_PORT=8080
STORK_DATABASE_HOST=/run/postgresql
STORK_DATABASE_NAME=stork
STORK_DATABASE_USER_NAME=stork
STORK_SERVER_ENABLE_METRICS=1
```
