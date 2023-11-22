FROM golang:1.20-alpine AS build-stage

WORKDIR /auth-service
ENV GOPATH=/

# Download packages only if module files changed
COPY services/auth/go.mod services/auth/go.sum ./
RUN go mod download

COPY /services/auth ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /auth /auth

EXPOSE 8040
EXPOSE 8041

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client
RUN chmod +x wait-4-postgres.sh

ENTRYPOINT ["/auth"]