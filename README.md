# Golang_rbtree

## 简介：
基于原始红黑树概念编写的红黑树map类,key和val都可以是interface{}类型，需要注意的是在创建红黑树时需要传入一个key值的比较方法CompareFunc

## 提供方法：
Map.Len() : 获取Map长度

Map.Add(key, val) : 向Map里添加一个键值对

Map.Delete(key) : 在Map里删除一个键值对

Map.Set(key, val) : 设置key的值为val, 需要确认key值存在，不然无法添加

Map.Get(key) : 获得key对应的val值，如果不存在返回会是false, nil, 存在会是true, val

Map.Range() : 获得Map对应键值对的Pair结构体信息通道chan，用于for range

## 举例：
在Test/rbtree_test.go文件中有测试代码

```go
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
			if a.(int) < b.(int) {
				return uint8(1)
			} else if a.(int) > b.(int) {
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
```
