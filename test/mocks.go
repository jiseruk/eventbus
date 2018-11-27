package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	lambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/client"
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
	environment.Variables["queue_name"] = aws.String(client.AWS_RESOURCE_PREFIX + "dead-letter-" + subscriber)
	environment.Variables["environment"] = config.GetCurrentEnvironment()

	createArgs := &lambda.CreateFunctionInput{
		Code:             createCode,
		FunctionName:     aws.String(client.AWS_RESOURCE_PREFIX + "lambda-" + subscriber),
		Handler:          aws.String("subscriber.handler"),
		Role:             aws.String(config.Get("engines.AWS.lambda.executionRole")),
		Runtime:          aws.String("python2.7"),
		Environment:      &environment,
		DeadLetterConfig: &lambda.DeadLetterConfig{TargetArn: &dlqArn},
		Tags:         map[string]*string{"project": aws.String("wequeue")},
	}
	return createArgs
}

func executeMockedRequest(router *gin.Engine, method string, uri string, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, strings.NewReader(body))
	router.ServeHTTP(rec, req)
	return rec
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
