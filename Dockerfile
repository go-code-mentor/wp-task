FROM golang:1.24.1 AS builder

COPY . /opt/app

WORKDIR /opt/app

RUN make tools
RUN make lint
RUN make test

RUN rm -rfv ./bin/*

RUN CGO_ENABLED=0 go build -o bin/work_planner ./cmd/app/main.go

FROM alpine:3.21

ARG TZ="Europe/Moscow"
ARG USER="work_planner"

ENV LANG='C.UTF-8'  \
    LC_ALL='C.UTF-8' \
    TZ=${TZ} \
    USER=${USER} \
    WORKDIR=/opt/app \
    PATH="/opt/app:${PATH}"

RUN apk add -U tzdata &&  \
    ln -fns /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo $TZ > /etc/timezone 

RUN adduser --shell /bin/bash --disabled-password --gecos "" ${USER}

COPY --from=builder --chown=${USER} /opt/app/bin ${WORKDIR}

USER ${USER}
WORKDIR ${WORKDIR}