FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash -u 1000 developer
