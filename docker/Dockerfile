FROM golang:1.21-bookworm AS builder

WORKDIR /build
COPY ["src/main.go", "test/polly.mp3", "/build/"]
RUN go build main.go

FROM alpine:3.19.0

WORKDIR /app
RUN apk update && apk add ffmpeg
COPY --from=builder /build/* /app

CMD [ "./main" ]
