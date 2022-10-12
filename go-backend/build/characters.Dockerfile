# Build application
FROM golang:1.17 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
RUN go build -o /out/characters ./cmd/characters/main.go

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/characters ./
EXPOSE 8081
CMD [ "./characters", "-port=8081", "-mode=release" ]