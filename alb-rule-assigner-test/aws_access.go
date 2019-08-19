package main

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var sess = session.Must(session.NewSession())
var cfnClient = cloudformation.New(sess)

type cfnInterface interface {
	CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	WaitUntilStackCreateComplete(input *cloudformation.DescribeStacksInput) error
	WaitUntilStackDeleteComplete(input *cloudformation.DescribeStacksInput) error
	ListStacks(input *cloudformation.ListStacksInput) (*cloudformation.ListStacksOutput, error)
}

type cfnManager interface {
	CreateStack(stackName *string, templatePath *string, parameters []*cloudformation.Parameter, testRunPrefix *string) error
	DeleteStack(stackName *string) error
	GetStacks(testRun *string) ([]*cloudformation.StackSummary, error)
}

type cfnManage struct {
	cfnClient        cfnInterface
	testStackTagName *string
}

var cmInstance = &cfnManage{
	cfnClient:        cfnClient,
	testStackTagName: aws.String("bc:internal:testrun"),
}

func (cm cfnManage) CreateStack(stackName *string, templatePath *string, parameters []*cloudformation.Parameter, testRunPrefix *string) error {

	// Read in the file as CreateStack take url or template string
	bb, err := ioutil.ReadFile(*templatePath)
	if err != nil {
		return err
	}

	// Create Stack
	csr := &cloudformation.CreateStackInput{
		TemplateBody: aws.String(string(bb)),
		StackName:    stackName,
		Capabilities: []*string{
			aws.String("CAPABILITY_AUTO_EXPAND"),
			aws.String("CAPABILITY_IAM"),
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		Parameters: parameters,
		Tags: []*cloudformation.Tag{
			&cloudformation.Tag{
				Key:   cm.testStackTagName,
				Value: testRunPrefix,
			},
		},
	}

	_, err = cm.cfnClient.CreateStack(csr)
	if err != nil {
		return err
	}

	// Wait Stack
	dsi := &cloudformation.DescribeStacksInput{
		StackName: stackName,
	}
	err = cm.cfnClient.WaitUntilStackCreateComplete(dsi)
	if err != nil {
		return err
	}

	// Success
	return nil
}

func (cm cfnManage) DeleteStack(stackName *string) error {

	// issue delete stack command
	dsr := &cloudformation.DeleteStackInput{
		StackName: stackName,
	}
	_, err := cm.cfnClient.DeleteStack(dsr)
	if err != nil {
		return nil
	}

	// wait delete stack command
	// Wait Stack
	dsi := &cloudformation.DescribeStacksInput{
		StackName: stackName,
	}
	err = cm.cfnClient.WaitUntilStackDeleteComplete(dsi)
	if err != nil {
		return err
	}
	return nil
}

func (cm cfnManage) GetStacks(testRun *string) ([]*cloudformation.StackSummary, error) {

	// get all stacks
	lsi := &cloudformation.ListStacksInput{
		StackStatusFilter: []*string{
			aws.String("CREATE_COMPLETE"),
			aws.String("CREATE_IN_PROGRESS"),
			aws.String("CREATE_FAILED"),
			aws.String("REVIEW_IN_PROGRESS"),
			aws.String("ROLLBACK_COMPLETE"),
			aws.String("ROLLBACK_IN_PROGRESS"),
			aws.String("UPDATE_COMPLETE"),
			aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"),
			aws.String("UPDATE_IN_PROGRESS"),
			aws.String("UPDATE_ROLLBACK_COMPLETE"),
			aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"),
			aws.String("UPDATE_ROLLBACK_FAILED"),
			aws.String("UPDATE_ROLLBACK_IN_PROGRESS"),
		},
	}
	lso, err := cfnClient.ListStacks(lsi)
	if err != nil {
		return nil, err
	}

	// filter for stacks of current test run
	filteredStack := make([]*cloudformation.StackSummary, 0)
	for _, stack := range lso.StackSummaries {
		if strings.HasPrefix(*stack.StackName, *testRun+"-") {
			filteredStack = append(filteredStack, stack)
		}
	}
	return filteredStack, nil
}
