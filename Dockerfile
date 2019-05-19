FROM golang:1.12-stretch

ENV GO111MODULE=on
WORKDIR /go/src/github.com/deelawn/BrainPaaswd

COPY . .

CMD go run main.go