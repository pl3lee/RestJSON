FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY . .
RUN ./scripts/buildprod.sh

FROM debian:stable-slim AS runner
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /app/restjson-api /usr/bin/restjson-api
CMD ["restjson-api"]
