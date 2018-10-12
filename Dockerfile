FROM golang:alpine
RUN apk add --no-cache curl 
RUN apk add --no-cache git 
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/wenance/wequeue-management_api
WORKDIR /go/src/github.com/wenance/wequeue-management_api
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . .
#RUN go build -o main . 
#CMD ["./main"]
EXPOSE 8080 
CMD ["go", "run", "main.go"]
