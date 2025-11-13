FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN apk add --no-cache make
RUN make

FROM alpine

WORKDIR /app
COPY --from=builder /app/bin/* .
RUN apk add --no-cache curl
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost || exit 1

EXPOSE 80

CMD [ "/app/prservice" ]
