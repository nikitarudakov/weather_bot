FROM golang:1.21.1

WORKDIR /bot

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /weatherBot

EXPOSE 8080

CMD ["/weatherBot"]