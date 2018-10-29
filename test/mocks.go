package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	lambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/model"
)

func getTopicMock(name string, engine string, resource string) *model.Topic {
	topic := model.Topic{ResourceID: resource, Name: name, Engine: engine}
	topic.CreatedAt = model.Clock.Now()
	topic.UpdatedAt = model.Clock.Now()
	return &topic
}

func getLambdaMock(endpoint string, subscriber string, topic string, dlqArn string) *lambda.CreateFunctionInput {
	contents, err := ioutil.ReadFile("/tmp/function.zip")
	if err != nil {
		return nil
	}

	createCode := &lambda.FunctionCode{
		ZipFile: contents,
	}

	environment := lambda.Environment{Variables: make(map[string]*string)}
	environment.Variables["subscriber_url"] = &endpoint
	environment.Variables["topic"] = &topic
	environment.Variables["queue_name"] = aws.String("dlq_lambda_" + subscriber)
	environment.Variables["environment"] = aws.String(config.LOCAL)

	createArgs := &lambda.CreateFunctionInput{
		Code:             createCode,
		FunctionName:     aws.String("lambda_subscriber_" + subscriber),
		Handler:          aws.String("subscriber.handler"),
		Role:             aws.String("arn:role:dummy"),
		Runtime:          aws.String("python2.7"),
		Environment:      &environment,
		DeadLetterConfig: &lambda.DeadLetterConfig{TargetArn: &dlqArn},
	}
	return createArgs
}

func executeMockedRequest(router *gin.Engine, method string, uri string, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, strings.NewReader(body))
	router.ServeHTTP(rec, req)
	return rec
}
