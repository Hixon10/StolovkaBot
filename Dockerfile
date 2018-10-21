FROM golang:1.11

COPY . /go/src/github.com/Hixon10/StolovkaBot
WORKDIR /go/src/github.com/Hixon10/StolovkaBot

RUN apt-get update && apt-get install -y

#RUN go get -u github.com/pions/webrtc

RUN go install -v ./...

CMD ["StolovkaBot"]