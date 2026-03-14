FROM golang:1.26-alpine3.23 AS build-stage
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o go-app .

FROM alpine:3.23 AS test-stage
COPY --from=build-stage /app/go-app /app/go-app
ENTRYPOINT [ "/app/go-app" ]
