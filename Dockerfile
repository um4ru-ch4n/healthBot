FROM golang:1.17

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bin/healthBot.app ./main.go 

CMD /bin/healthBot.app