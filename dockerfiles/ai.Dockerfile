FROM golang:1.20-alpine AS build-stage

WORKDIR /ai-svc
ENV GOPATH=/

RUN mkdir -p src/services
RUN mkdir ai

# Download packages only if module files changed
COPY go.mod go.sum ./
RUN go mod download

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait-4-postgres.sh

COPY /gen ./gen
COPY /src/services/ai ./src/services/ai
COPY /src/internal ./src/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /ai src/services/ai/cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /ai /ai

EXPOSE 8020

ENTRYPOINT ["/ai"]