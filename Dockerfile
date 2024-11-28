FROM golang:1.23.2-alpine3.20 AS builder
COPY *.go /src/
COPY go.* /src/
WORKDIR /src
RUN  go build .
RUN ls -la /src/

FROM alpine:3.20
COPY --from=builder /src/notes-telegram /opt/app/
WORKDIR /opt/app
ENTRYPOINT [ "/opt/app/notes-telegram" ]

