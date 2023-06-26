FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main/

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /myapp
COPY --from=builder /build/server ./server
COPY cert.crt key.key ./
ENV PORT_REST=443
ENV PORT_gRPC=5051
EXPOSE $PORT_REST
EXPOSE $PORT_gRPC
ENTRYPOINT ["./server"]
# command arg - --pCertFile  --pKeyFile

