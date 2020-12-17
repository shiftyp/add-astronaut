FROM golang:1.12-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o bin/launch

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/launch /usr/local/bin/
CMD ["launch"]
