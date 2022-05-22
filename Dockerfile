FROM golang:1.14-alpine AS builder
 
WORKDIR /srv/go-app
ADD . .
RUN go build -o arctic-fox


FROM debian:buster
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/config.json .
COPY --from=builder /srv/go-app/arctic-fox .

CMD ["./arctic-fox"]