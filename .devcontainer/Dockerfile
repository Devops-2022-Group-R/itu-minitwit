# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.202.1/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN apt-get update \
    && apt-get install -y sqlite3 libsqlite3-dev

RUN git config --global core.autocrlf true
