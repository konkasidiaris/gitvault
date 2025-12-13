FROM mcr.microsoft.com/devcontainers/go:1.25-trixie

RUN apt-get update && apt-get install -y ca-certificates openssl git curl && \
    rm -rf /var/lib/apt/lists/*

USER vscode
