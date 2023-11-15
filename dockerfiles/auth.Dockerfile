FROM golang:1.20-alpine AS build-stage

WORKDIR /auth-svc
ENV GOPATH=/

# Download packages only if module files changed
COPY /services/auth/go.mod /services/auth/go.sum ./
RUN go mod download

COPY /services/auth ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /auth /auth

EXPOSE 8001
EXPOSE 8010

ENTRYPOINT ["/auth"]