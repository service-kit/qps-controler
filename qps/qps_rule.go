package qps

import (
	"container/list"
	"sync"
	"time"
)

type RuleFunc func(interface{}) bool

type QPSFresher struct {
	interval        int64
	lastRefreshTime int64
}

func (qf *QPSFresher) refresh(now int64) bool {
	if now > qf.lastRefreshTime+qf.interval {
		qf.lastRefreshTime = now
		return true
	}
	return false
}

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

func (c *QPSCounter) OverLimit() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count >= c.limit
}

type QPSRuleItem struct {
	ruleFunc RuleFunc
	counter  QPSCounter
	fresher  QPSFresher
}

func (qri *QPSRuleItem) SetLimit(limit int64) {
	qri.counter.SetLimit(limit)
}

func (qri *QPSRuleItem) Pass(cond interface{}) bool {
	if nil != qri.ruleFunc {
		isInRule := qri.ruleFunc(cond)
		if !isInRule {
			return false
		}
	}
	return qri.counter.Add(1)
}

func (qri *QPSRuleItem) fresh(now int64) {
	if qri.fresher.refresh(now) {
		qri.counter.Reset()
	}
}

type QPSRule struct {
	ruleList *list.List
	mutex    sync.RWMutex
}

func (qr *QPSRule) Init() {
	qr.ruleList = new(list.List)
}

func (qr *QPSRule) addRule(rf RuleFunc, limit int64, duration time.Duration) {
	qr.mutex.Lock()
	defer qr.mutex.Unlock()
	if 0 == duration {
		duration = time.Second
	}
	qr.ruleList.PushBack(
		&QPSRuleItem{
			ruleFunc: rf,
			counter: QPSCounter{
				limit: limit,
			},
			fresher: QPSFresher{
				interval:        int64(duration),
				lastRefreshTime: time.Now().UnixNano(),
			},
		},
	)
}

func (qr *QPSRule) pass(cond interface{}) bool {
	qr.mutex.Lock()
	defer qr.mutex.Unlock()
	for i := qr.ruleList.Front(); nil != i; i = i.Next() {
		if !i.Value.(*QPSRuleItem).Pass(cond) {
			return false
		}
	}
	return true
}

func (qr *QPSRule) refresh(now int64) bool {
	qr.mutex.Lock()
	defer qr.mutex.Unlock()
	for i := qr.ruleList.Front(); nil != i; i = i.Next() {
		i.Value.(*QPSRuleItem).fresh(now)
	}
	return true
}
