FROM golang:1.19

WORKDIR /bot

COPY . ./

RUN go build -o bin/skele main.go

CMD ["./bin/skele"]
