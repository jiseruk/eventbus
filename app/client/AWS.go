package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

const SNSEndpoint = "http://localstack:4575"

type AWSEngine struct {
	SNSClient snsiface.SNSAPI
}

func GetSNSClient() snsiface.SNSAPI {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		//Credentials: credentials.NewSharedCredentials("", "default"),
		Credentials: credentials.NewStaticCredentials("foo", "bar", ""),
		Endpoint:    aws.String(SNSEndpoint),
	})
	if err != nil {

	}
	svc := sns.New(sess)
	return svc
}

//CreateTopic lala
func (aw *AWSEngine) CreateTopic(name string) (*CreateTopicOutput, error) {
	// Make svc.AddPermission request
	var input = &sns.CreateTopicInput{Name: &name}
	snsoutput, err := aw.SNSClient.CreateTopic(input)
	if err != nil {
		fmt.Printf("Error: %#v", err)
		return nil, err
	}
	output := &CreateTopicOutput{Resource: *snsoutput.TopicArn}
	return output, nil
}

func (aw AWSEngine) GetName() string {
	return "AWS"
}

