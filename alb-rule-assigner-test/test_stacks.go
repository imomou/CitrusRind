package main

import (
	"math/rand"
	"time"
	"strings"
	"os"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type testStackManager interface {
	CreateStack(stackName *string, templatePath *string, parameters []*cloudformation.Parameter) error
	DeleteStack(stackName *string) error
}

type testStackManage struct {
	cm cfnManager
	stackParams []*cloudformation.Parameter
	randomPrefix string
}

var stmInstance = &testStackManage{
	cm: cmInstance,
	stackParams: make([]*cloudformation.Parameter, 0),
	randomPrefix: "XVlBzg", //RandStringRunes(10),
}

func (tsm *testStackManage) SuiteSetupTeardown(s *godog.Suite) error {

	lambdaZipUrl := os.Getenv("LAMBDA_ZIP_URL")

	s.BeforeSuite(func() {

	})
	s.AfterSuite(func() {
		stacks, err := tsm.cm.GetStacks(&tsm.randomPrefix)
		if err != nil {
			panic(err)
		}

		for _, stack := range stacks {
			tsm.cm.DeleteStack(stack.StackName)
		}
	})
	
	s.BeforeFeature(func(feature *gherkin.Feature) {
		for _, step := range feature.Background.Steps {
			step.Text = strings.Replace(step.Text, "<prefix>", tsm.randomPrefix, -1)
			step.Text = strings.Replace(step.Text, "<lambdas3path>", lambdaZipUrl, -1)
		}
	})
	s.BeforeScenario(func(scenarioObject interface{}) {
		scenario, ok := scenarioObject.(*gherkin.Scenario)
		if !ok {
			panic("wasn't a scenario passed in.")
		}
		
		for _, step := range scenario.Steps {
			step.Text = strings.Replace(step.Text, "<prefix>", tsm.randomPrefix, -1)
			step.Text = strings.Replace(step.Text, "<lambdas3path>", lambdaZipUrl, -1)
		}
		
	})

	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
