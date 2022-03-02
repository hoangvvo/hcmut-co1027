FROM golang:1.17

WORKDIR /usr/src/app

RUN apt-get update
RUN apt-get install dos2unix

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./main.go

CMD ["app"]