FROM golang:1.25-trixie AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

WORKDIR /app/examples/helloworld

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -v -ldflags="-X main.gitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -o /app/helloworld ./...

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/helloworld .

EXPOSE 9000

CMD ["/app/helloworld"]
