FROM golang:1.20-alpine AS build

WORKDIR /go/src/github.com/nicolaka/tfbi

ARG tag="v0.1"

COPY . .

RUN go build \
    -ldflags="-X main.Version=${tag} -X main.BuildDate=$(date '+%Y%m%d-%H:%M:%S')"

FROM alpine:3.19 AS prod

COPY --from=build /go/src/github.com/nicolaka/tfbi/tfbi /bin/tfbi

USER nobody

ENTRYPOINT [ "/bin/tfbi" ]
