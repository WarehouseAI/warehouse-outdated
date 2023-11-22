FROM golang:1.20-alpine AS build-stage

WORKDIR /ai-service
ENV GOPATH=/

# Download packages only if module files changed
COPY services/ai/go.mod services/ai/go.sum ./
RUN go mod download

COPY /services/ai ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /ai /ai

EXPOSE 8020
EXPOSE 8021

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client
RUN chmod +x wait-4-postgres.sh

ENTRYPOINT ["/ai"]