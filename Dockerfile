FROM golang:1.20-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add ca-certificates

WORKDIR /
COPY services.upload services.upload
COPY services.shared ../services.shared
COPY keys keys 
WORKDIR /services.upload
ENV CGO_ENABLED=0
COPY ./services.upload/go.mod ./services.upload/go.sum ./
RUN  --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . . 
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main ./src/cmd/main.go

FROM scratch

ENV PORT 8080

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /services.upload/main .
COPY --from=builder /keys ./keys
COPY --from=builder /services.upload/src/locales ./src/locales

EXPOSE $PORT

CMD ["/main"]