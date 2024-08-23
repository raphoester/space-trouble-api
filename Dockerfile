FROM golang:1.22-alpine3.19 AS build


ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /root
COPY . .

ENV GOCACHE=/root/.cache/go-build
ENV GOPATH=/root/go
RUN --mount=type=cache,target="/root/.cache/go-build" \
	--mount=type=cache,target="/root/go/pkg/mod" \
    go build -o /root/bin/api ./cmd/api;


# runner final image
FROM alpine:3 AS cmd

ARG SYSTEM_USER=space-trouble-api

RUN addgroup --gid 1000 -S ${SYSTEM_USER} && adduser --uid 1000 -S ${SYSTEM_USER} -G ${SYSTEM_USER}

USER ${SYSTEM_USER}

WORKDIR /home/app/${SYSTEM_USER}

COPY --chown=${SYSTEM_USER}:${SYSTEM_USER} --from=build /root/bin/ ./

ENTRYPOINT "./api"