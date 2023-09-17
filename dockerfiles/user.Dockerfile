FROM golang:1.20-alpine AS build-stage

WORKDIR /user-svc
ENV GOPATH=/

RUN mkdir -p src/services
RUN mkdir user

# Download packages only if module files changed
COPY go.mod go.sum ./
RUN go mod download

# Download alpine package and install psql-client for the script
COPY wait-4-postgres.sh ./
RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait-4-postgres.sh

COPY /gen ./gen
COPY /src/services/user ./src/services/user
COPY /src/internal ./src/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /user src/services/user/cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /user /user

EXPOSE 8000
EXPOSE 8001

ENTRYPOINT ["/user"]