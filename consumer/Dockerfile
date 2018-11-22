FROM golang:latest
RUN mkdir -p /go/src/github.com/wenance/wequeue-management_api/consumer
WORKDIR /go/src/github.com/wenance/wequeue-management_api/consumer
COPY . .
RUN ls -ltr
EXPOSE 9000 
CMD ["go", "run", "consumer.go"]