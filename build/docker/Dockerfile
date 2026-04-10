FROM openapitools/openapi-generator-cli:latest AS api

WORKDIR /local

COPY api api
COPY scripts scripts
COPY .openapi-generator-ignore .openapi-generator-ignore

RUN docker-entrypoint.sh generate -c /local/scripts/generator-cfg.yaml

ENTRYPOINT ["bash"]

############################################

FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ internal/
COPY cmd/ cmd/
COPY --from=api /local/ ./

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux \
      go build \
      -ldflags="-w -s" \
      -installsuffix 'static' \
      -o ./hospital-webapi-srv ./cmd/hospital-api-service

############################################

FROM scratch

LABEL org.opencontainers.image.authors="Filip Pauco"
LABEL org.opencontainers.image.title="Hospital Ward WebAPI Service"
LABEL org.opencontainers.image.description="WEBAPI for managing hospital ward visits and beds"

COPY --from=build /app/hospital-webapi-srv ./

EXPOSE 8080
ENTRYPOINT ["./hospital-webapi-srv"]
