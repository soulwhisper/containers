# Containers

## About

This repo contains a collection of container images updated automatically to keep up with upstream versions. The images try to adhere to the following principles:

- Images are built upon a [Alpine](https://hub.docker.com/_/alpine) or [Ubuntu](https://hub.docker.com/_/ubuntu) base image.
- No use of [s6-overlay](https://github.com/just-containers/s6-overlay).
- Semantic versioning is available to specify exact versions to run.
- The container filesystem must be able to be immutable.

## Available Images

Images can be browsed on the GitHub Packages page for this repo's [packages](https://github.com/soulwhisper?tab=packages&repo_name=containers).
