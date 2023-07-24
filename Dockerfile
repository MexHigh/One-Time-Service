ARG BUILD_FROM
FROM $BUILD_FROM

RUN apk add --no-cache caddy
COPY Caddyfile /etc/caddy/Caddyfile

CMD [ "/usr/sbin/caddy", "run", "--config", "/etc/caddy/Caddyfile" ]