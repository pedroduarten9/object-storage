#
# Build stage
#
FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /api cmd/api.go

#
# Release stage
#
FROM docker

WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 3000

ENTRYPOINT ["/api"]