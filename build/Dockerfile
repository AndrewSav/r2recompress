ARG DEBIAN_FRONTEND=noninteractive
FROM golang:1.21.0-bullseye
RUN apt-get update && apt-get install -y curl gnupg apt-transport-https && \
  curl -sS https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | tee /etc/apt/trusted.gpg.d/microsoft.gpg > /dev/null && \
  sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-debian-bullseye-prod bullseye main" > /etc/apt/sources.list.d/microsoft.list' && \
  apt-get update && apt-get install -y powershell
