FROM --platform=${BUILDPLATFORM} golang:1.21 as build

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

COPY go.mod go.sum /app/
COPY cmd/ /app/cmd/
COPY vendor/ /app/vendor/
COPY internal/ /app/internal/

ARG GIT_TAG
ARG GIT_COMMIT
ARG GIT_BRANCH

RUN BUILD_TIME=$(date -Iseconds) \
    go build \
    -o rsstodiscord \
    -ldflags "\
        -s -w \
        -X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Version=${GIT_TAG}' \
        -X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Commit=${GIT_COMMIT}' \
        -X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Branch=${GIT_BRANCH}' \
        -X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Timestamp=${BUILD_TIME}' \
    " \
    cmd/rsstodiscord/main.go

FROM --platform=${BUILDPLATFORM} alpine:3.12 AS security

RUN apk add --no-cache \
    ca-certificates

ENV USER=rsstodiscord
ENV UID=1000
ENV GID=1000

RUN addgroup \
        -g "$GID" \
        -S \
        $USER \
    && adduser \
        -S \
        -D \
        -g "" \
        -G "$USER" \
        -H \
        -u "$UID" \
        "$USER"

FROM --platform=${TARGETPLATFORM} scratch

EXPOSE 9090

ARG GIT_TAG
ARG GIT_COMMIT
ARG GIT_BRANCH

LABEL version=${GIT_TAG}
LABEL build.branch=${GIT_BRANCH}
LABEL build.sha=${GIT_COMMIT}

COPY LICENSE /app/
COPY --from=security /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=security /etc/passwd /etc/passwd
COPY --from=build /app/rsstodiscord /app/

EXPOSE 8080

ENTRYPOINT [ "/app/rsstodiscord" ]

CMD [ "serve" ]
