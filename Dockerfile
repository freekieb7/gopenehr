FROM golang:1.23-alpine

#RUN apk update && apk add --no-cache build-base cmake git

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN go mod download

CMD ["air", "-c", ".air.toml"]