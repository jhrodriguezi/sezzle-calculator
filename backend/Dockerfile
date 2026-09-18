FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
