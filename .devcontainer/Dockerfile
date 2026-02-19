FROM golang:1.26.0-trixie

ARG DOCKER_GID
ARG USERNAME=gitvault
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt update && apt upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates openssl git curl && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian trixie stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-ce-cli && \
    groupadd --gid $USER_GID $USERNAME && \
    useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    groupadd -f -g ${DOCKER_GID} docker && \
    usermod -aG docker ${USERNAME} && \
    apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

USER ${USERNAME}
