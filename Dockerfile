FROM golang:latest as builder
#RUN apk add --no-cache curl 
#RUN apk add --no-cache git 
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/wenance/wequeue-management_api
WORKDIR /go/src/github.com/wenance/wequeue-management_api
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
#For local environment
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x ./wait-for-it.sh
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .
RUN swag init
RUN cd test && go test -covermode=count -coverprofile=cover.out -coverpkg=../app/...
#RUN bin/tests.sh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 

FROM busybox:musl
WORKDIR /go/src/github.com/wenance/wequeue-management_api 
COPY --from=builder /go/src/github.com/wenance/wequeue-management_api/main /go/src/github.com/wenance/wequeue-management_api/main 
COPY --from=builder /go/src/github.com/wenance/wequeue-management_api/lambda /go/src/github.com/wenance/wequeue-management_api/lambda
COPY --from=builder /go/src/github.com/wenance/wequeue-management_api/app/config/local.yml /go/src/github.com/wenance/wequeue-management_api/app/config/local.yml 

CMD ["./main"]
EXPOSE 8080 
#CMD ["go", "run", "main.go"]
