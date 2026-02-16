FROM golang:1.24-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH:-amd64}" go build -trimpath -ldflags="-s -w" -o /out/dmh-api ./api/dmh.go


FROM alpine:3.20

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=build /out/dmh-api /app/dmh-api
COPY api/etc /app/etc

EXPOSE 8080

ENTRYPOINT ["/app/dmh-api"]
CMD ["-f", "/app/etc/dmh-api.prod.yaml"]
