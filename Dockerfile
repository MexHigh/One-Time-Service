ARG BUILD_FROM

FROM golang:1.18 AS builder
WORKDIR /go/src/app
COPY backend/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' .

FROM $BUILD_FROM
WORKDIR /app
# copy compiled backend
COPY --from=builder /go/bin/backend /app/backend
# copy frontends
COPY frontend-internal/ /frontend-internal
COPY frontend-public/ /frontend-public

CMD [ "/app/backend" ]