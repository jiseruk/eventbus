FROM golang:latest as builder
#RUN apk add --no-cache curl 
#RUN apk add --no-cache git 
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/wenance/wequeue-management_api
WORKDIR /go/src/github.com/wenance/wequeue-management_api
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only -v
#For local environment
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x ./wait-for-it.sh
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .

FROM builder as tests
WORKDIR /go/src/github.com/wenance/wequeue-management_api
RUN swag init
RUN cd test && go test -covermode=count -coverprofile=cover.out -coverpkg=../app/...
#RUN bin/tests.sh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 

#FROM busybox:musl
FROM alpine:latest
RUN apk add --no-cache curl bash jq
ENV GOPATH /go
WORKDIR /go/src/github.com/wenance/wequeue-management_api 
COPY bin/entrypoint-vault.sh /entrypoint/
COPY --from=tests /etc/ssl/certs /etc/ssl/certs
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/main /go/src/github.com/wenance/wequeue-management_api/main 
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/lambda /go/src/github.com/wenance/wequeue-management_api/lambda
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/app/config/local.yml /go/src/github.com/wenance/wequeue-management_api/app/config/local.yml 
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/app/config/develop.yml /go/src/github.com/wenance/wequeue-management_api/app/config/develop.yml 
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/app/config/stage.yml /go/src/github.com/wenance/wequeue-management_api/app/config/stage.yml 
COPY --from=tests /go/src/github.com/wenance/wequeue-management_api/app/config/prod.yml /go/src/github.com/wenance/wequeue-management_api/app/config/prod.yml 
ENTRYPOINT ["/entrypoint/entrypoint-vault.sh"]
CMD ["./main"]
EXPOSE 8080 
#CMD ["go", "run", "main.go"]
