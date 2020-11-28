FROM golang:1-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on
RUN go build -o server ./example/...

FROM alpine:3

COPY --from=builder /app/server /server
COPY --from=builder /app/example/frontend /frontend
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/server"]
