Feature: basic auto assign alb priority
  To save time looking up ALB rule priority
  As a software engineer
  I want to rule to be automatically allocated without my input
  
  Note tokens with arrow bracket such as <token> are replaced at
  run time with before feature func.
  
  Background:
    Given A Lambda Zip path of <lambdas3path>
    When I provision alb_rule_assigner_macro.template with stack name <prefix>-alb-rule-assigner-macro
    Then A cloudformation stack <prefix>-alb-rule-assigner-macro should exist

  Scenario: Assign priority from empty alb listener
    When I provision single_rule_assigned.template with stack name <prefix>-single-rule
    Then A cloudformation stack <prefix>-single-rule should exist

#  Scenario: Assign priority from empty alb listener
#    When I provision single_rule_assigned.template with stack name <prefix>-single-rule
#    Then A cloudformation stack <prefix>-alb-rule-assigner-macro should exist
