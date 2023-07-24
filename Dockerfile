FROM golang:1.18-alpine AS builder
WORKDIR /go/src/app
COPY backend/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' .

FROM alpine:latest
WORKDIR /app
# copy compiled backend
COPY --from=builder /go/bin/backend /app/backend
# copy frontends
COPY frontend-internal/ /frontend-internal
COPY frontend-public/ /frontend-public

LABEL \
  io.hass.version="0.1.0" \
  io.hass.type="addon" \
  io.hass.arch="amd64"

CMD [ "/app/backend" ]