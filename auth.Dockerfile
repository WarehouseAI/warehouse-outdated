FROM golang:1.20-alpine AS build-stage

WORKDIR /auth-svc
ENV GOPATH=/

RUN mkdir -p src/services
RUN mkdir auth

# Download packages only if module files changed
COPY go.mod go.sum ./
RUN go mod download

COPY /gen ./gen
COPY /src/services/auth ./src/services/auth
COPY /src/internal ./src/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth src/services/auth/cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /auth /auth

EXPOSE 8010

ENTRYPOINT ["/auth"]