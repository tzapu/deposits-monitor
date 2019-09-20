FROM golang:1.12

WORKDIR /app

ADD go.mod go.mod
ADD go.sum go.sum
RUN go mod download

ADD . .
RUN go build

ENTRYPOINT ["./deposits-monitor"]