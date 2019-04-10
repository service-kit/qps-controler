package qps

import (
	"errors"
)

type Rule struct {
	ruleFunc RuleFunc
	limit int64
}

type QPSControler struct {
	ruleMap map[string]*QPSRule
}

func (qc *QPSControler) Init() {
	qc.ruleMap = make(map[string]*QPSRule)
}

func (qc *QPSControler) NewRule(ruleName string,rules ...Rule) error {
	if nil != qc.ruleMap[ruleName] {
		return errors.New("rule exist")
	}
	newRule := new(QPSRule)
	newRule.Init()
	for _,rule := range rules {
		newRule.AddRule(rule.ruleFunc,rule.limit)
	}
	qc.ruleMap[ruleName] = newRule
	return nil
}