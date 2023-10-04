FROM golang:1.20-alpine AS backend-builder
WORKDIR /go/src/app
COPY backend/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' .

FROM node:20 AS internal-frontend-builder
WORKDIR /tmp
COPY frontend-internal/ .
RUN yarn install --network-timeout=600000
RUN yarn run build

FROM alpine:latest
# install stuff
RUN apk add --no-cache tzdata
# copy in backend and frontends 
WORKDIR /app
COPY --from=backend-builder /go/bin/backend /app/backend
COPY --from=internal-frontend-builder /tmp/build /frontend-internal
COPY frontend-public/ /frontend-public

ARG ADDON_VERSION=dev
LABEL \
  io.hass.version=$ADDON_VERSION \
  io.hass.type="addon"

CMD [ "/app/backend", "-db", "/share/one-time-service/db.json" ]