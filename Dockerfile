## Build
FROM golang:1.18-buster AS build

WORKDIR /

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /app

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app /app
COPY /ui ./ui
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]