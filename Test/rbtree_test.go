package Test

import (
	"fmt"
	"rbtree/rbmap"
	"testing"
)

func TestMapAddFunc(t *testing.T) {
	mp := rbmap.NewMap(func(a, b interface{}) uint8 {
		if a.(int) < b.(int) {
			return uint8(1)
		} else if a.(int) > b.(int) {
			return uint8(2)
		} else {
			return uint8(0)
		}
	})
	for i := 1; i <= 10; i++ {
		mp.Add(i, i+1)
	}
	//fmt.Println(mp.Len())
	for i := 1; i <= 9; i++ {
		mp.Delete(i)
	}
	for pair := range mp.Range() {
		fmt.Println(pair.Key, pair.Val)
	}
}
