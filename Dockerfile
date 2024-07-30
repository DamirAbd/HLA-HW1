# syntax=docker/dockerfile:1

FROM golang:1.22.1 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY . .

  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

FROM scratch AS build-release-stage
  WORKDIR /

  COPY --from=build-stage /api /api

  EXPOSE 8080

  ENTRYPOINT ["/api"]