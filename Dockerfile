FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /myapp ./cmd

FROM alpine:latest AS run

COPY --from=build /myapp /myapp

WORKDIR /app
EXPOSE 8000
CMD ["/myapp"]
