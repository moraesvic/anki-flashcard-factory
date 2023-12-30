FROM golang:1.21-bookworm AS builder

WORKDIR /build
COPY src /build
RUN go build

FROM debian:12.4-slim

WORKDIR /app
RUN apt-get update && apt-get install -y ffmpeg
COPY --from=builder /build/* /app/
RUN cp flashcard-factory /bin/

CMD [ "bash" ]
