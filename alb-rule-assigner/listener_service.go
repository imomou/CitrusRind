package main

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

// GeneratorService Constructor Poco
type GeneratorService struct {
	elbv2Service *elbv2.ELBV2
}

func newGeneratorService(elbv2Service *elbv2.ELBV2) *GeneratorService {
	return &GeneratorService{elbv2Service: elbv2Service}
}

// GetListernerRules Get Listener Rules
func (service *GeneratorService) GetListernerRules(ListenerArn string) ([]*elbv2.Rule, error) {

	results, err := service.elbv2Service.DescribeRules(&elbv2.DescribeRulesInput{
		ListenerArn: &ListenerArn})

	return results.Rules, err
}

// GetRandomRules : Expected <=listenRulecap, -1 not able to randomly generate rule
func (service *GeneratorService) GetRandomRules(rules []*elbv2.Rule, listenRulecap int) (int, error) {

	iter := 0
	//https://flaviocopes.com/go-random/
	//https://www.calhoun.io/creating-random-strings-in-go/
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

// ReplaceFragment Inside the map, replace Priority key
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
