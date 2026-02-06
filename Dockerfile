FROM busybox AS bin
COPY ./dist /binaries
RUN if [[ "$(arch)" == "x86_64" ]]; then \
        architecture="amd64"; \
    else \
        architecture="arm64"; \
    fi; \
    cp /binaries/ch-linux-${architecture} /bin/ch && \
    chmod +x /bin/ch && \
    chown 1000:1000 /bin/ch

FROM chainguard/wolfi-base
LABEL org.opencontainers.image.title="ContainerHive"
LABEL org.opencontainers.image.description="Swarm it. Build it. Run it."
LABEL org.opencontainers.image.ref.name="main"
LABEL org.opencontainers.image.licenses='GPLv3'
LABEL org.opencontainers.image.vendor="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.authors="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.url="https://github.com/timo-reymann/ContainerHive"
LABEL org.opencontainers.image.documentation="https://github.com/timo-reymann/ContainerHive"
LABEL org.opencontainers.image.source="https://github.com/timo-reymann/ContainerHive.git"
RUN apk add --no-cache bash \
    && adduser -D -u 1000 container-hive

ARG BUILD_TIME
ARG BUILD_VERSION
ARG BUILD_COMMIT_REF
LABEL org.opencontainers.image.created=$BUILD_TIME
LABEL org.opencontainers.image.version=$BUILD_VERSION
LABEL org.opencontainers.image.revision=$BUILD_COMMIT_REF

COPY --from=bin /bin/ch /bin/ch
WORKDIR /workspace
ENTRYPOINT [ "/bin/ch" ]
