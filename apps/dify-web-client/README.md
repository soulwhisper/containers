# dify-web-client

Containerised [Dify webapp-conversation](https://github.com/langgenius/webapp-conversation) — the official Next.js web client template for [Dify](https://dify.ai) chatbot apps.

## Quick start

```sh
docker run -p 3000:3000 \
  -e NEXT_PUBLIC_APP_ID=<your-app-id> \
  -e NEXT_PUBLIC_APP_KEY=<your-app-key> \
  -e NEXT_PUBLIC_API_URL=https://api.dify.ai/v1 \
  ghcr.io/soulwhisper/dify-web-client:latest
```

Then open http://localhost:3000.

## Configuration

All configuration is done via `NEXT_PUBLIC_*` environment variables at build time (Next.js inlines these into the client bundle).

| Variable | Description |
|---|---|
| `NEXT_PUBLIC_APP_ID` | Dify app ID (from the app URL) |
| `NEXT_PUBLIC_APP_KEY` | API key from the app's "API Access" page |
| `NEXT_PUBLIC_API_URL` | Dify API base URL (default: `https://api.dify.ai/v1`) |

To customise the app title, description, or locale, fork the upstream template and build your own image.

## Building

```sh
just local-build dify-web-client
```

This builds a local image tagged `dify-web-client:<short-sha>` and runs container tests against it.
