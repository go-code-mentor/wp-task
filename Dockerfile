FROM golang:1.24.1 AS builder

COPY . /opt/app

WORKDIR /opt/app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v2.0.2

RUN ./bin/golangci-lint run
RUN go test ./... -v
RUN rm -rfv ./bin/*
RUN CGO_ENABLED=0 go build -o bin/work_planner ./cmd/app/main.go

FROM golang:1.24.1-alpine

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