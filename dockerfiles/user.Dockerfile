FROM golang:1.20-alpine AS build-stage

WORKDIR /user-service
ENV GOPATH=/

# Download packages only if module files changed
COPY services/user/go.mod services/user/go.sum ./
RUN go mod download

COPY /services/user ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /user cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /user /user

EXPOSE 8000
EXPOSE 8001

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client
RUN chmod +x wait-4-postgres.sh

ENTRYPOINT [ "/user" ]