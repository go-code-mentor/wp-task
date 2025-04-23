FROM golang:1.24.1 AS builder

COPY . /opt/app

WORKDIR /opt/app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./lint v2.0.2
RUN ./lint/golangci-lint run
RUN go test ./... -v
RUN CGO_ENABLED=0 go build -o bin/work_planner ./cmd/app/main.go

FROM alpine:3.21

ARG USER="work_planner"

ENV USER=${USER} \
    WORKDIR=/opt/app \
    PATH="/opt/app:${PATH}"

RUN adduser --shell /bin/bash --disabled-password --gecos "" ${USER}

COPY --from=builder --chown=${USER} /opt/app/bin ${WORKDIR}

USER ${USER}
WORKDIR ${WORKDIR}