FROM golang:1.20-alpine AS build-stage

WORKDIR /stat-service
ENV GOPATH=/

# Download packages only if module files changed
COPY services/stat/go.mod services/stat/go.sum ./
RUN go mod download

COPY /services/stat ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /stat cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /stat /stat

EXPOSE 8022
EXPOSE 8023

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client
RUN chmod +x wait-4-postgres.sh

ENTRYPOINT ["/stat"]