package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

// GeneratorService blah blah
type GeneratorService struct {
	elbv2Service *elbv2.ELBV2
}

func newGeneratorService(elbv2Service *elbv2.ELBV2) *GeneratorService {
	return &GeneratorService{elbv2Service: elbv2Service}
}

// GetListernerRules Get Listener Rules
func (service *GeneratorService) GetListernerRules(ListenerArn string) []*elbv2.Rule {

	results, err := service.elbv2Service.DescribeRules(&elbv2.DescribeRulesInput{
		ListenerArn: &ListenerArn})

	if err != nil {
		fmt.Println(err)
	}

	return results.Rules
}

// GetRandomRules stuff
func (service *GeneratorService) GetRandomRules(rules []*elbv2.Rule, listenRulecap int) (int, error) {

	iter := 0
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	l := listenRulecap + 1

	for iter < math.MaxInt32 {
		n := r.Intn(l)

		if !containPriority(rules, n) && n > 0 {
			return n, nil
		}
		iter++
	}

	return -1, errors.New("Unable to allocate errors")
}

// ReplaceFragment stuff
func (service *GeneratorService) ReplaceFragment(properties map[string]interface{}, value int) map[string]interface{} {

	properties["Priority"] = value

	return properties
}

func containPriority(rules []*elbv2.Rule, generatedRule int) bool {

	for _, v := range rules {

		if *v.Priority == strconv.Itoa(generatedRule) {
			return true
		}
	}
	return false
}
