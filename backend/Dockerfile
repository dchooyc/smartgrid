# Build stage: compiles the server. The tests also run in this stage.
FROM golang:1.27 AS build
WORKDIR /src

# Dependencies first, so they are only downloaded again when they change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO is off so the binary needs nothing from the image it runs in.
RUN CGO_ENABLED=0 go build -o /server .

# Run stage: a small image holding only the server.
FROM alpine:3
COPY --from=build /server /server
USER nobody
EXPOSE 8080
ENTRYPOINT ["/server"]
