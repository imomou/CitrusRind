package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func Test_GetRandomRules(t *testing.T) {

	elbvRules := []*elbv2.Rule{
		{Priority: aws.String("1")},
		{Priority: aws.String("2")}}

	generatorService := newGeneratorService(nil)

	rand, err := generatorService.GetRandomRules(elbvRules, 3)

	fmt.Println(rand)

	if rand != 3 {
		t.Errorf("failed to generate rule %s", err)
	}
}

func Test_ReplaceFragment(t *testing.T) {

	properties := map[string]interface{}{
		"Priority": "321"}

	generatorService := newGeneratorService(nil)

	properties = generatorService.ReplaceFragment(properties, 123)

	if properties["Priority"] != 123 {
		t.Errorf("failed to assign priority %s", properties["Priority"])
	}
}

func Test_ReplaceFragment_WithoutPriority(t *testing.T) {

	properties := map[string]interface{}{}

	generatorService := newGeneratorService(nil)

	properties = generatorService.ReplaceFragment(properties, 123)

	if properties["Priority"] != 123 {
		t.Errorf("failed to assign priority %s", properties["Priority"])
	}
}
