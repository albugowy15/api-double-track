FROM golang:1.22.4 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-double-track cmd/api/main.go

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/static-debian12 AS build-release-stage
WORKDIR /
COPY --from=build-stage /api-double-track /api-double-track
COPY app.env.example app.env
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api-double-track"]

LABEL org.opencontainers.image.source=https://github.com/albugowy15/api-double-track
LABEL org.opencontainers.image.description="Go-based REST API that offers a variety of endpoints designed to supply data for integration with the Double Track Recommendation Web"
LABEL org.opencontainers.image.licenses=MIT
