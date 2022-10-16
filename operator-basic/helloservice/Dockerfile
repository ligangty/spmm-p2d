# Build the manager binary
FROM quay.io/fedora/fedora:35

USER root

RUN dnf install -y glibc

ADD build/hello /usr/bin/hello

RUN chmod +x /usr/bin/hello

ENV HOST 0.0.0.0
ENV PORT 8080
# ENV GIN_MODE release
EXPOSE 8080


ENTRYPOINT ["/usr/bin/hello"]
