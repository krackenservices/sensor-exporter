FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git build-base lm-sensors-dev

WORKDIR /app

COPY sensor-exporter/ ./sensor-exporter/

WORKDIR /app/sensor-exporter

RUN go mod tidy
RUN go build -o /sensor-exporter

RUN ls

FROM alpine:3.19

RUN apk add --no-cache lm-sensors lm-sensors-sensord lm-sensors-detect

COPY --from=builder /sensor-exporter /usr/local/bin/sensor-exporter

EXPOSE 9255

ENTRYPOINT ["/usr/local/bin/sensor-exporter"]
CMD ["--web.listen-address=:9255"]

