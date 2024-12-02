FROM golang:1.23-alpine  AS build

WORKDIR /go/src/github.com/nicolaka/tfbi

ARG tag="v0.2"

COPY . .

RUN go build \
    -ldflags="-X main.Version=${tag} -X main.BuildDate=$(date '+%Y%m%d-%H:%M:%S')"

FROM alpine:3.20 AS prod

COPY --from=build /go/src/github.com/nicolaka/tfbi/tfbi /bin/tfbi

USER nobody

ENTRYPOINT [ "/bin/tfbi" ]
