# Build application
FROM golang:1.17 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY build ./
RUN go mod download
RUN go build -o /out/updater ./cmd/updater/main.go

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/updater ./
EXPOSE 8180
CMD [ "./updater", "-port=8180", "-mode=release" ]