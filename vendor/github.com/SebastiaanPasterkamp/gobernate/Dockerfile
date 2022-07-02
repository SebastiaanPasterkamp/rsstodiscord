# See https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
FROM golang:1.15-stretch AS builder

# Git is used for dependencies
RUN apt-get update \
    && apt-get install -y \
        git \
        ca-certificates

ENV USER=gobernate
ENV UID=1000
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    -u "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/gobernate/

COPY go.mod go.sum ./

RUN apt-get install -y \
      build-essential \
    && go mod download \
    && go mod verify

COPY . .

ARG PROJECT="gobernate"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /go/bin/gobernate \
    -ldflags="\
        -w -s \
        -X ${PROJECT}/version.Name=${PROJECT} \
        -X ${PROJECT}/version.Release=$(git describe --tags --always) \
        -X ${PROJECT}/version.Commit=$(git rev-parse --verify HEAD) \
        -X ${PROJECT}/version.BuildTime=$(date -u -Iseconds)" \
    cmd/main.go

FROM scratch

COPY --from=builder \
    /etc/passwd \
    /etc/group \
    /etc/
COPY --from=builder \
    /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/
COPY --from=builder \
    /go/bin/gobernate \
    /go/bin/gobernate

USER 1000

ENV PORT 8080

ENTRYPOINT [ "/go/bin/gobernate" ]
