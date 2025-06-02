# See https://aka.ms/devcontainer for more information.
# This is a basic devcontainer for Go development.

FROM mcr.microsoft.com/devcontainers/go:1-1.24

# [Optional] Install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Install common Go tools.
RUN go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# [Optional] Set the default shell to bash rather than sh.
ENV SHELL /bin/bash

# [Optional] Set up a non-root user for VS Code.
# ARG USERNAME=vscode
# ARG USER_UID=1000
# ARG USER_GID=$USER_UID
# RUN groupadd --gid $USER_GID $USERNAME \
#     && useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME \
#     && apt-get update && apt-get install -y sudo \
#     && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
#     && chmod 0440 /etc/sudoers.d/$USERNAME

# [Optional] Expose ports.
EXPOSE 8080
