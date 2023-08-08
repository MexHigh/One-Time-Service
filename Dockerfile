FROM golang:1.18-alpine AS backend-builder
WORKDIR /go/src/app
COPY backend/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' .

FROM node:20 AS internal-frontend-builder
WORKDIR /tmp
COPY frontend-internal/ .
RUN npm install && npm run build

FROM alpine:latest
WORKDIR /app
# copy compiled backend
COPY --from=backend-builder /go/bin/backend /app/backend
# copy frontends
COPY --from=internal-frontend-builder /tmp/build /frontend-internal
COPY frontend-public/ /frontend-public

LABEL \
  io.hass.version="0.1.0" \
  io.hass.type="addon" \
  io.hass.arch="amd64"

CMD [ "/app/backend", "-db", "/share/one-time-service/db.json" ]