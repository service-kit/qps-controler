package qps

import (
	"container/list"
	"sync"
)

type RuleFunc func(interface{}) bool

type QPSCounter struct {
	mutex sync.RWMutex
	count int64
	limit int64
}

func (c *QPSCounter) SetLimit(limit int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.limit = limit
}

func (c *QPSCounter) Add(v int64) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.count+v > c.limit {
		return false
	}
	c.count += v
	return true
}

func (c *QPSCounter) Value() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *QPSCounter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.count = 0
}

type QPSRuleItem struct {
	ruleFunc RuleFunc
	counter  QPSCounter
}

func (qri *QPSRuleItem) SetLimit(limit int64) {
	qri.counter.SetLimit(limit)
}

func (qri *QPSRuleItem) Pass(cond interface{}) bool {
	isInRule := qri.ruleFunc(cond)
	if isInRule {
		return qri.counter.Add(1)
	}
	return true
}

type QPSRule struct {
	ruleList *list.List
	mutex    sync.RWMutex
}

func (qr *QPSRule) Init() {
	qr.ruleList = new(list.List)
}

func (qr *QPSRule) AddRule(rf RuleFunc,limit int64) {
	qr.mutex.Lock()
	defer qr.mutex.Unlock()
	qr.ruleList.PushBack(&QPSRuleItem{ruleFunc: rf,counter:QPSCounter{limit:limit}})
}

func (qr *QPSRule) Pass(cond interface{}) bool {
	qr.mutex.Lock()
	defer qr.mutex.Unlock()
	for i := qr.ruleList.Front(); nil != i; i = i.Next() {
		if !i.Value.(*QPSRuleItem).Pass(cond) {
			return false
		}
	}
	return true
}
