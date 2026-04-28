FROM golang:1.25-alpine AS build
RUN apk add --no-cache build-base ca-certificates git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ordersystem ./cmd/ordersystem

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/ordersystem ./ordersystem
COPY migrations ./migrations
COPY configs ./configs
ENV MIGRATIONS_DIR=/app/migrations
EXPOSE 8080 50051 8081
ENTRYPOINT ["/app/ordersystem"]
