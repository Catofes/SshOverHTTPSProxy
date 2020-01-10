FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download
RUN env CGO_ENABLED=0 go build -o /SshOverHTTPS .

FROM alpine:3.6
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /main /usr/bin/SshOverHTTPS
ENTRYPOINT ["/usr/bin/SshOverHTTPS"]