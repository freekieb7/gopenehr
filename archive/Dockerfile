FROM golang:1.24.3 AS development

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /gopenehr cmd/main.go

CMD ["/gopenehr"]

FROM scratch AS production
WORKDIR /
COPY --from=development /gopenehr /
CMD ["/gopenehr"]