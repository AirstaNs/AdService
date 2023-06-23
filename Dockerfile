FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main/

FROM alpine
#RUN #apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build ./app
EXPOSE 8080
CMD ["./app/server"]


