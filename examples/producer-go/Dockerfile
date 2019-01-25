FROM golang:1.10

WORKDIR /go/src/producer-go
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v ./...

CMD ["producer-go"]