FROM golang:1.20-alpine AS dev

RUN apk add --no-cache \
    g++ \
    inotify-tools

WORKDIR /go/src/github.com/kaizendorks/terraform-cloud-exporter

ENTRYPOINT ["./hot-reload.sh"]

FROM golang:1.20-alpine AS build

WORKDIR /go/src/github.com/kaizendorks/terraform-cloud-exporter

ARG tag="v0.0.0"
ARG sha="hash_commit"

COPY . .
RUN go build \
    -ldflags="-X main.Version=${tag} -X main.Commit=${sha} -X main.BuildDate=$(date '+%Y%m%d-%H:%M:%S')"

FROM alpine:3.18 AS prod

COPY --from=build /go/src/github.com/kaizendorks/terraform-cloud-exporter/terraform-cloud-exporter /bin/terraform-cloud-exporter

USER nobody
ENTRYPOINT [ "/bin/terraform-cloud-exporter" ]
