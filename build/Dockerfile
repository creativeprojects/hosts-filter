FROM alpine:latest

LABEL maintainer Fred <fred@gcreativeprojects.tech>

COPY hosts-filter /usr/bin/

VOLUME /hosts-filter
WORKDIR /hosts-filter

ENTRYPOINT ["hosts-filter"]
CMD ["--help"]
