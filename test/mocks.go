package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/model"
)

func getTopicMock(name string, engine string, resource string, owner string, description string) *model.Topic {
	topic := model.Topic{ResourceID: resource, Name: name, Engine: engine, Owner: owner, Description: description}
	topic.CreatedAt = model.Clock.Now()
	topic.UpdatedAt = model.Clock.Now()
	topic.SecurityToken = UUIDMock{}.GetUUID()
	return &topic
}

func getSubscriberMock(name string, topic string, Type string, resource string) *model.Subscriber {
	subscriber := model.Subscriber{Name: name, Topic: topic, Type: Type, ResourceID: resource}
	subscriber.CreatedAt = model.Clock.Now()
	subscriber.UpdatedAt = model.Clock.Now()
	return &subscriber
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
	environment.Variables["queue_name"] = aws.String(client.GetAWSResourcePrefix() + "dead-letter-" + subscriber)
	environment.Variables["environment"] = config.GetCurrentEnvironment()

	createArgs := &lambda.CreateFunctionInput{
		Code:             createCode,
		FunctionName:     aws.String(client.GetAWSResourcePrefix() + "lambda-" + subscriber),
		Handler:          aws.String("subscriber.handler"),
		Role:             aws.String(config.Get("engines.AWS.lambda.executionRole")),
		Runtime:          aws.String("python2.7"),
		Environment:      &environment,
		DeadLetterConfig: &lambda.DeadLetterConfig{TargetArn: &dlqArn},
		Tags:             map[string]*string{"project": aws.String("wequeue")},
		VpcConfig: &lambda.VpcConfig{
			SecurityGroupIds: []*string{aws.String(config.Get("engines.AWS.lambda.securityGroupId"))},
			SubnetIds:        config.GetArray("engines.AWS.lambda.subnetIds"),
		},
	}
	return createArgs
}

func executeMockedRequest(router *gin.Engine, method string, uri string, body string, headers ...string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, strings.NewReader(body))
	for _, h := range headers {
		header := strings.Split(h, ":")
		req.Header[header[0]] = []string{header[1]}
	}
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

type UUIDMock struct{}

func (u UUIDMock) GetUUID() string {
	return "uuid"
}
