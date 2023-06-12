package Test

import (
	"fmt"
	"rbtree/rbmap"
	"testing"
)

func TestMapAddFunc(t *testing.T) {
	// 获取一个map, 设置比较方法
	mp := rbmap.NewMap(
		func(a, b interface{}) uint8 {
			c, d := a.(int), b.(int)
			if c < d {
				return uint8(1)
			} else if c > d {
				return uint8(2)
			} else {
				return uint8(0)
			}
		})
	// 往map里存10个键值对
	for i := 1; i <= 10; i++ {
		mp.Add(i, i+1)
	}
	// 输出map的长度
	// expected 10
	fmt.Println(mp.Len())

	// 删除1--9的键
	for i := 1; i <= 9; i++ {
		mp.Delete(i)
	}
	// 获取key=2的值输出（可以发现不存在）
	_, val := mp.Get(2)
	// expected <nil>
	fmt.Println(val)

	// 通过Range方法循环整个树
	for pair := range mp.Range() {
		// expected 10 11 because of only one key
		fmt.Println(pair.Key, pair.Val)
	}
}
