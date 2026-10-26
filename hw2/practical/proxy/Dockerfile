FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/proxy .

FROM alpine:3.21

RUN adduser -D -H -u 65532 proxy
USER proxy
COPY --from=build /out/proxy /usr/local/bin/proxy
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/proxy"]
