package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// CfnFragment Fragment
type CfnFragment struct {
	Properties map[string]interface{} `json:"Properties,omitempty"`
	Type       interface{}            `json:"Type,omitempty"`
	Condition  string                 `json:"Condition,omitempty"`
}

// LambdaCFNRequest structure of request from CFN macro
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html
type LambdaCFNRequest struct {
	Region         string                 `json:"region"`
	AccountID      string                 `json:"accountId"`
	Fragment       CfnFragment            `json:"fragment"`
	TransformID    string                 `json:"transformId"`
	Params         map[string]interface{} `json:"params"`
	RequestID      string                 `json:"requestId"`
	TemplateParams map[string]interface{} `json:"templateParameterValues"`
}

// LambdaCFNResponse structure of response from CFN macro
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html
type LambdaCFNResponse struct {
	RequestID string      `json:"requestId"`
	Status    string      `json:"status"`
	Fragment  interface{} `json:"fragment"`
}

// HandlerRequest handling of lambda request
func HandlerRequest(request LambdaCFNRequest) (*LambdaCFNResponse, error) {

	listenRulecap, err := strconv.Atoi(os.Getenv("RULECAP"))
	region := os.Getenv("AWS_REGION")
	params := request.Params

	ListenerArn := fmt.Sprintf("%v", params["target"])

	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})

	if err != nil {
		log.Fatal(fmt.Sprintf("Bad request %s", err))
	}

	client := elbv2.New(sess)
	service := newGeneratorService(client)
	cfnResponse := &LambdaCFNResponse{RequestID: request.RequestID}

	listeners := service.GetListernerRules(ListenerArn)
	randRule, err := service.GetRandomRules(listeners, listenRulecap)

	if err != nil {

		log.Fatal(fmt.Sprintf("%s", err))
		cfnResponse.Status = "failure"

	} else {

		cfnResponse.Status = "success"
	}

	fragmentResources := request.Fragment
	fragmentResources.Properties = service.ReplaceFragment(fragmentResources.Properties, randRule)

	if err != nil {

		log.Fatal(fmt.Sprintf("Bad request %s", err))
	}

	cfnResponse.Fragment = fragmentResources
	return cfnResponse, err
}

func main() {

	lambda.Start(HandlerRequest)
}
