FROM golang:1.19

WORKDIR /bot

COPY . ./

RUN go build -o bin/skele main.go

CMD ["./bin/skele"]

# make it listen to master and pull changes && restart when master to pushed
