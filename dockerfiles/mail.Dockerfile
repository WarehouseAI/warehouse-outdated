FROM golang:1.20-alpine AS build-stage

WORKDIR /mail-service
ENV GOPATH=/

# Download packages only if module files changed
COPY services/mail/go.mod services/mail/go.sum ./
RUN go mod download

COPY /services/mail ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /mail cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

EXPOSE 5672

COPY --from=build-stage /mail /mail

ENTRYPOINT ["/mail"]