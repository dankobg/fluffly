ARG GO_VERSION=1.24-alpine
FROM golang:${GO_VERSION} AS dev
RUN apk update && apk upgrade \
  && apk add --no-cache git \
  && go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "main.go", "serve"]