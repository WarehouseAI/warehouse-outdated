FROM golang:1.20-alpine AS build-stage

WORKDIR /auth-svc
ENV GOPATH=/

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /auth /auth

RUN ls

EXPOSE 8080

ENTRYPOINT ["/auth"]