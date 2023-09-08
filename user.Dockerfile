FROM golang:1.20-alpine AS build-stage

WORKDIR /user-svc
ENV GOPATH=/

RUN mkdir -p src/services
RUN mkdir user

# Download packages only if module files changed
COPY go.mod go.sum ./
RUN go mod download

COPY /gen ./
COPY /src/services/user ./src/services/

RUN CGO_ENABLED=0 GOOS=linux go build src/services/user/cmd/main.go -o /user

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /user /user

EXPOSE 8000
EXPOSE 8001

CMD ["/user"]