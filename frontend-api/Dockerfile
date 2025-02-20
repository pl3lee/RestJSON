FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY . .
RUN ./scripts/buildprod.sh

FROM debian:stable-slim AS runner
COPY --from=builder /app/webjson-frontend-api /usr/bin/webjson-frontend-api
CMD ["webjson-frontend-api"]
