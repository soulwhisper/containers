// dsh-web-relay: bind the container's concrete address and forward to
// 127.0.0.1:<DSH_UPSTREAM_PORT> where dsh listens. dsh >= 0.2.1-alpha.2
// refuses wildcard binds; this relay keeps the Service/LAN path without
// patching dsh. Design borrowed from runzhliu/deepseek-harness-docker (MIT).
//
// Env:
//   DSH_PORT                 relay listen port (default 3080)
//   DSH_UPSTREAM_PORT        dsh's loopback port (default 3080; dsh owns it)
//   DSH_WEB_RELAY_INTERFACE  interface name to bind when the host has
//                            multiple non-loopback IPv4 addresses
"use strict";

const net = require("net");
const os = require("os");

const port = Number(process.env.DSH_PORT || 3080);
const upstream = {
  host: "127.0.0.1",
  port: Number(process.env.DSH_UPSTREAM_PORT || 3080),
};

function pickAddress() {
  const want = process.env.DSH_WEB_RELAY_INTERFACE;
  const ifaces = os.networkInterfaces();
  const candidates = [];
  for (const [name, addrs] of Object.entries(ifaces)) {
    for (const a of addrs || []) {
      const family = typeof a.family === "string" ? a.family : `IPv${a.family}`;
      if (family !== "IPv4" || a.internal) continue;
      if (want && name !== want) continue;
      candidates.push({ name, address: a.address });
    }
  }
  if (want) {
    if (candidates.length === 0) {
      console.error(`relay: DSH_WEB_RELAY_INTERFACE=${want} has no non-loopback IPv4`);
      process.exit(1);
    }
    return candidates[0];
  }
  if (candidates.length === 0) return { name: "lo", address: "127.0.0.1" };
  if (candidates.length > 1) {
    console.error(
      `relay: multiple non-loopback IPv4 addresses (${candidates
        .map((c) => `${c.name}=${c.address}`)
        .join(", ")}); set DSH_WEB_RELAY_INTERFACE`
    );
    process.exit(1);
  }
  return candidates[0];
}

const bind = pickAddress();

const server = net.createServer((client) => {
  const up = net.connect(upstream, () => client.pipe(up).pipe(client));
  up.on("error", (err) => {
    console.error(`relay: upstream error: ${err.message}`);
    client.destroy();
  });
  client.on("error", () => up.destroy());
});

server.on("error", (err) => {
  console.error(`relay: cannot bind ${bind.address}:${port}: ${err.message}`);
  process.exit(1);
});

server.listen(port, bind.address, () => {
  console.log(
    `relay: ${bind.address}:${port} -> ${upstream.host}:${upstream.port} (iface ${bind.name})`
  );
});
