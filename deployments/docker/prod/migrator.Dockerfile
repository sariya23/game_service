FROM ubuntu:20.04 AS migrations
WORKDIR /app
ENV PATH="/usr/lib/postgresql/16/bin:${PATH}"
RUN apt-get update && \
    apt-get install -y postgresql-client git make
ADD --chmod=755 https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
COPY migrations ./migrations
COPY Makefile ./
COPY config ./config
