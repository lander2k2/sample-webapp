FROM golang:1.13 as builder
ENV GO111MODULE=on
WORKDIR /build
COPY main.go .
RUN CGO_ENABLED=0 go build -o app main.go

FROM alpine:3.11
COPY --from=builder /build/app /app
COPY index.html /index.html
CMD ["/app"]

