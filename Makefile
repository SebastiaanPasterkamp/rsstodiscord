BINARY_NAME="rsstodiscord"
PWD=$(shell pwd)
UID=$(shell id -u)
GID=$(shell id -g)
GIT_TAG=$(shell git describe master --tags 2> /dev/null || echo -n "v1.0.0")
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date)

update:
	go get -u all
	go mod tidy -compat=1.17
	go mod vendor

coverage:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out
	go tool cover -func cover.out

lint:
	golangci-lint run -v

format:
	gofmt -s -w cmd
	gofmt -s -w internal

test:
	go test -race -count 10 -v ./...

build:
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOFLAGS=-mod=vendor \
	go build \
		-o $(BINARY_NAME) \
		-a \
		-ldflags "\
			-X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Version=$(GIT_TAG)' \
			-X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Commit=$(GIT_COMMIT)' \
			-X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Branch=$(GIT_BRANCH)' \
			-X 'github.com/SebastiaanPasterkamp/rsstodiscord/internal/build.Timestamp=$(BUILD_TIME)' \
		" \
		cmd/$(BINARY_NAME)/main.go

docker:
	docker build \
		--build-arg GIT_BRANCH="$(GIT_BRANCH)" \
		--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
		--build-arg GIT_TAG="$(GIT_TAG)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--tag $(BINARY_NAME) \
		.

clean:
	rm \
		$(BINARY_NAME) \
		cover.out
