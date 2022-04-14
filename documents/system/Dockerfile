FROM caddy:2.4.5-builder AS builder

RUN xcaddy build \
    --with github.com/greenpau/caddy-auth-portal \
    --with github.com/greenpau/caddy-auth-jwt

FROM caddy:2.4.5

RUN apk update
RUN apk add --no-cache nss-tools curl

RUN echo -e "{\n admin 0.0.0.0:2019\n }" > /etc/caddy/Caddyfile

COPY --from=builder /usr/bin/caddy /usr/bin/caddy