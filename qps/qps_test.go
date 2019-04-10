package qps

import (
	"fmt"
	"testing"
	"time"
)

func Test_QPS(t *testing.T) {
	qc := QPSControler{}
	qc.Init()
	qc.NewRule("test",
		Rule{
			RuleFunc: func(i interface{}) bool {
				if nil == i {
					return false
				}
				if val, ok := i.(string); !ok || "" == val {
					return false
				}
				return true
			},
			Limit:    10,
			Duration: time.Second,
		},
		Rule{
			RuleFunc: func(i interface{}) bool {
				if nil == i {
					return false
				}
				if val, ok := i.(int); !ok || 0 == val {
					return false
				}
				return true
			},
			Limit:    10,
			Duration: time.Second,
		},
	)
	for {
		for i := 0; i < 12; i++ {
			fmt.Print(qc.Pass("test", nil), " ")
		}
		fmt.Println()
		for i := 0; i < 12; i++ {
			fmt.Print(qc.Pass("test", ""), " ")
		}
		fmt.Println()
		for i := 0; i < 12; i++ {
			fmt.Print(qc.Pass("test", "1"), " ")
		}
		fmt.Println()
		for i := 0; i < 12; i++ {
			fmt.Print(qc.Pass("test", 0), " ")
		}
		fmt.Println()
		for i := 0; i < 12; i++ {
			fmt.Print(qc.Pass("test", 1), " ")
		}
		fmt.Println()
		fmt.Println()
		fmt.Println()
		time.Sleep(750 * time.Millisecond)
	}
}
