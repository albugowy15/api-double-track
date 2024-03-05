FROM golang:1.22.0 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /api-double-track cmd/server/main.go

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:latest AS run-stage
RUN adduser -D -u 1001 appuser
WORKDIR /app
RUN chown -R appuser: /app
RUN apk --no-cache add ca-certificates
COPY --from=build-stage /api-double-track /app/api-double-track
COPY  app.env.example /app/app.env
EXPOSE 8080
USER appuser
ENTRYPOINT [ "/app/api-double-track" ]
