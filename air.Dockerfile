FROM golang:alpine

WORKDIR /app

RUN apk --update --no-cache add make curl
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download

HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost || exit 1
CMD ["air", "-c", ".air.toml"]
