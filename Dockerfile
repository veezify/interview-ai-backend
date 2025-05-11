FROM golang:latest

WORKDIR /app

COPY . ./
RUN go mod download

COPY *.go ./

RUN go build -o /go-docker-demo

EXPOSE 8080

CMD [ "/go-docker-demo" ]
