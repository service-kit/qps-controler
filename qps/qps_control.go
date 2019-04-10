package qps

import (
	"errors"
	"time"
)

type Rule struct {
	RuleFunc RuleFunc
	Limit    int64
	Duration time.Duration
}

type QPSControler struct {
	ruleMap    map[string]*QPSRule
	freshTimer *time.Ticker
}

func (qc *QPSControler) Init() {
	qc.ruleMap = make(map[string]*QPSRule)
	qc.freshTimer = time.NewTicker(10 * time.Millisecond)
	go qc.fresh()
}

func (qc *QPSControler) NewRule(ruleName string, rules ...Rule) error {
	if nil != qc.ruleMap[ruleName] {
		return errors.New("rule exist")
	}
	newRule := new(QPSRule)
	newRule.Init()
	for _, rule := range rules {
		newRule.addRule(rule.RuleFunc, rule.Limit, rule.Duration)
	}
	qc.ruleMap[ruleName] = newRule
	return nil
}

func (qc *QPSControler) Pass(ruleName string, cond interface{}) bool {
	if nil == qc.ruleMap[ruleName] {
		return true
	}
	return qc.ruleMap[ruleName].pass(cond)
}

func (qc *QPSControler) fresh() {
	for {
		select {
		case <-qc.freshTimer.C:
			now := time.Now().UnixNano()
			for _, rule := range qc.ruleMap {
				rule.refresh(now)
			}
			qc.freshTimer = time.NewTicker(10 * time.Millisecond)
		}
	}
}
