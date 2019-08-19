package main

import (
	//"fmt"
	"regexp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	
	"github.com/DATA-DOG/godog"
)

func (tsm *testStackManage) aLambdaZipPathOfLambdasPath(lambdaZipUrl string) error {
    // caveat, only bucketnames with characters and dashes works.
	r := regexp.MustCompile(`^https://(?P<Region>[\w\-]*)\.amazonaws\.com/(?P<BucketName>[\w\-]*)/(?P<Key>.*)$`)

	// fmt.Printf("%#v\n", r.FindStringSubmatch(`https://ap-southeast-1.amazonaws.com/bucket-name/this/is/a/test.zip`))
	// []string{"https://ap-southeast-1.amazonaws.com/bucket-name/this/is/a/test.zip", "ap-southeast-1", "bucket-name", "this/is/a/test.zip"}
	// fmt.Printf("%#v\n", r.SubexpNames())
	// []string{"", "Region", "BucketName", "Key"}
	values := r.FindStringSubmatch(lambdaZipUrl)

	tsm.stackParams = append(tsm.stackParams, &cloudformation.Parameter{
		ParameterKey:   aws.String("LambdaBucket"),
		ParameterValue: &values[2],
	}, &cloudformation.Parameter{
		ParameterKey:   aws.String("LambdaPath"),
		ParameterValue: &values[3],
	})

	return nil
}

func (tsm *testStackManage) iProvisionListener_ruletemplate(templatePath string, stackName string) error {
	err := tsm.cm.CreateStack(&stackName, aws.String("./templates/" + templatePath), tsm.stackParams, &tsm.randomPrefix)
	if err != nil {
		return err
	}

	// cleanup the parameters
	tsm.stackParams = make([]*cloudformation.Parameter, 0)
	return nil
}

func aCloudformationStackStackname(stackName string) error {
    return nil
}

func FeatureContext(s *godog.Suite) error {

	s.Step(`^A Lambda Zip path of (.*)$`, stmInstance.aLambdaZipPathOfLambdasPath)
    s.Step(`^I provision (.*) with stack name (.*)$`, stmInstance.iProvisionListener_ruletemplate)
    s.Step(`^A cloudformation stack (.*) should exist$`, aCloudformationStackStackname)
    
	return stmInstance.SuiteSetupTeardown(s)
}
