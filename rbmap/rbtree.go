// Package rbmap: 对红黑树的实现
package rbmap

import (
	"errors"
)

// golang 不支持重载运算符，所以只能通过方法调用
// 由于红黑树性质，所以树高不会太高，所以操作统一采用递归处理

// 红黑树五大定义：
// 1.节点只有红色和黑色
// 2.根节点是黑色
// 3.叶子节点是黑色
// 4.红色节点的儿子必是黑色节点
// 5.从任意节点出发，到达任意一个叶子节点的路径上的黑色节点数相同

// 由于定义4，所以从根到叶子的最长的可能路径不会多于最短的可能路径的两倍长。

// CompareFunc 定义比较方法 0：a==b, 1: a < b, 2 : a > b
type CompareFunc func(a, b interface{}) uint8

var (
	// ErrNodeAlreadyExists 插入节点时节点早已存在时报错
	ErrNodeAlreadyExists = errors.New("node already exists")
	ErrNodeNotExists     = errors.New("node not exists")
)

// Map 自定义的Map,提供常用接口
type Map struct {
	root *Node
	size int
	// 0：a==b, 1: a < b, 2 : a > b
	compareFunc CompareFunc
}

// NewMap 传入比较key值的函数作为构造方法
func NewMap(compareFunc CompareFunc) *Map {
	return &Map{
		root:        newLeaf(),
		size:        0,
		compareFunc: compareFunc,
	}
}

// Pair 键值对结构体
type Pair struct {
	Key keyItem
	Val valItem
}

// public:

// Range 根据中序遍历的方式循环红黑树，返回对应键值对
func (m *Map) Range() <-chan Pair {
	ch := make(chan Pair)
	go func() {
		m.ran(m.root, ch)
		close(ch)
	}()
	return ch
}

// Len 获得树的节点个数（存放值的节点个数）
func (m *Map) Len() int {
	return m.size
}

// Add 添加节点 key, val 如果节点存在就设置val的值并且返回错误
func (m *Map) Add(key keyItem, val valItem) error {
	node := m.findNode(m.root, key)
	if node.isLeaf() {
		m.insertNode(node, key, val)
		m.size++
		return nil
	}
	node.val = val
	return ErrNodeAlreadyExists
}

// Delete 根据key值删除对应节点， 如果节点不存在返回错误
func (m *Map) Delete(key keyItem) error {
	node := m.findNode(m.root, key)
	if node.isLeaf() {
		return ErrNodeNotExists
	}
	m.eraseNode(node)
	return nil
}

// Set 设置节点 key的值为val, 如果节点key不存在就返回false, 存在就修改返回true
func (m *Map) Set(key keyItem, val valItem) bool {
	node := m.findNode(m.root, key)
	if node.isLeaf() {
		return false
	}
	node.val = val
	return true
}

// Get 通过键值key找到对应的val,如果没有返回false
func (m *Map) Get(key keyItem) (bool, valItem) {
	node := m.findNode(m.root, key)
	if node.isLeaf() {
		return false, nil
	}
	return true, node.val
}

// 输出树结构，测试用
//func (m *Map) Print() {
//	m.print(m.root)
//}

// private:

//func (m *Map) print(node *Node) {
//	if !node.isLeaf() {
//		m.print(node.left)
//		fmt.Print("key: ", node.key, " color: ", node.color, " isLeft?: ", !node.isRoot() && node.isLeft(), " parent: ")
//		if !node.isRoot() {
//			fmt.Println(node.parent.key)
//		} else {
//			fmt.Println("None")
//		}
//		m.print(node.right)
//	}
//}

// 寻找节点，因为红黑树树高不会太高，所以选择递归寻找
func (m *Map) findNode(node *Node, key keyItem) *Node {
	if node.isLeaf() {
		return node
	}
	c := m.compareFunc(key, node.key)
	if c == 1 {
		return m.findNode(node.left, key)
	} else if c == 2 {
		return m.findNode(node.right, key)
	} else {
		return node
	}
}

// 插入节点，并且设置key和val
func (m *Map) insertNode(node *Node, key keyItem, val valItem) {
	node.key = key
	node.val = val
	node.color = RED
	node.left = newLeaf()
	node.right = newLeaf()
	node.left.parent = node
	node.right.parent = node
	// 进行插入调整
	m.insertSort(node)
}

// 对插入节点进行调整
func (m *Map) insertSort(node *Node) {
	if node.isRoot() {
		// 如果是根节点就设置成黑色（定义1）即可
		m.root = node
		node.color = BLACK
	} else if node.parent.isRed() {
		// 若是父节点颜色是黑色，就不需要处理当前节点，如果是红色就进行分类讨论
		if node.getUncle().isRed() {
			// 如果叔父节点颜色也是红色，就将当前父节点和叔父节点颜色变成黑色，祖父节点颜色变成红色然后对祖父节点进行调整
			node.getUncle().color, node.parent.color = BLACK, BLACK
			node.getGrandParent().color = RED
			m.insertSort(node.getGrandParent())
		} else {
			// 如果叔父节点是黑色，当前情况就是父节点和当前节点是红色，祖父节点是黑色，叔父节点是黑色，进行分类讨论
			// 获取当前节点是父节点的左儿子还是右儿子以及父节点是祖父节点的左儿子还是右儿子
			isLeft, isParentLeft := node.isLeft(), node.parent.isLeft()
			// 分四种情况讨论，分别是左左， 左右， 右左， 右右
			grandParent := node.getGrandParent()
			if isLeft && isParentLeft {
				// 左左， 同方向交换父节点以及祖父节点颜色然后右旋祖父节点
				node.parent.color, grandParent.color = BLACK, RED
				m.rotateRight(grandParent)
			} else if isLeft && !isParentLeft {
				// 左右，交换当前节点与祖父节点颜色，后对父节点右旋，祖父节点左旋
				node.color, grandParent.color = BLACK, RED
				m.rotateRight(node.parent)
				m.rotateLeft(grandParent)
			} else if !isLeft && isParentLeft {
				// 右左，交换当前节点与祖父节点颜色，然后对父节点左旋，祖父节点右旋
				node.color, grandParent.color = BLACK, RED
				m.rotateLeft(node.parent)
				m.rotateRight(grandParent)
			} else {
				// 右右，交换父节点与祖父节点颜色，然后对祖父节点左旋
				node.parent.color, grandParent.color = BLACK, RED
				m.rotateLeft(grandParent)
			}
		}
	}
}

// 删除节点node
func (m *Map) eraseNode(node *Node) {
	if node.left.isLeaf() {
		// 如果节点没有子节点或者只有右子节点，就用右子节点代替当前节点
		rightChild := node.right
		parent := node.parent
		rightChild.parent = parent
		if parent != nil {
			if node.isLeft() {
				parent.left = rightChild
			} else {
				parent.right = rightChild
			}
		}
		// 只有删除黑色节点才需要调整
		if node.isBlack() {
			m.eraseSort(rightChild)
		}
		node = nil
	} else if node.right.isLeaf() {
		// 如果左节点非空并且右节点是空节点
		leftChild := node.left
		parent := node.parent
		leftChild.parent = parent
		if parent != nil {
			if node.isLeft() {
				parent.left = leftChild
			} else {
				parent.right = leftChild
			}
		}
		// 同理只有删除黑色节点才需要调整
		if node.isBlack() {
			m.eraseSort(leftChild)
		}
		node = nil
	} else {
		// 如果节点有左右子节点，就找前继节点进行替换再删除前继节点
		leftMostChild := m.getLeftMostChild(node)
		node.key, node.val = leftMostChild.key, leftMostChild.val
		m.eraseNode(leftMostChild)
	}
}

// 寻找对应node节点的前继节点
func (m *Map) getLeftMostChild(node *Node) *Node {
	leftChild := node.left
	for !leftChild.isLeaf() {
		leftChild = leftChild.right
	}
	return leftChild.parent
}

// 删除节点后的调整
func (m *Map) eraseSort(node *Node) {
	// 调整时分类讨论
	if node.isRoot() {
		// 如果调整节点是根节点，设置成黑色并且更新根节点
		node.color = BLACK
		m.root = node
	} else if node.isRed() {
		// 如果调整节点是红色，就设置成黑色即可
		node.color = BLACK
	} else {
		// 如果当前节点是黑色，那么就分类讨论兄弟节点的颜色
		sibling := node.getSibling()
		if sibling.isRed() {
			// 如果兄弟节点是红色，那么设置兄弟节点是黑色，设置父节点是红色，并且调整父节点进行旋转（使得新兄弟节点变成黑色），再调整当前节点
			sibling.color = BLACK
			node.parent.color = RED
			// 通过旋转把父节点变成自己兄弟节点，此时更新后的兄弟节点就必是黑色了
			if node.isLeft() {
				m.rotateLeft(node.parent)
			} else {
				m.rotateRight(node.parent)
			}
			m.eraseSort(node)
		} else {
			// 如果兄弟节点是黑色，继续分类讨论树的形状以及兄弟节点的儿子节点的两个节点颜色
			if sibling.left.isBlack() && sibling.right.isBlack() {
				// 如果兄弟节点的两个子节点颜色都为黑色，直接将兄弟节点的颜色改为红色，再去调整父节点
				sibling.color = RED
				m.eraseSort(node.parent)
			} else {
				// 否则两个子节点肯定有一个节点颜色是红色，那么就可以进行之后操作，之后看树的形状以及对应颜色
				// 变换是使得与自己对称的节点的颜色为红色再进行旋转，如果自己是左子节点，那么就得让兄弟节点的右子节点变成红色，反之亦然
				if node.isLeft() && sibling.right.isBlack() {
					sibling.color, sibling.left.color = RED, BLACK
					m.rotateRight(sibling)
					sibling = node.getSibling()
				}
				if node.isRight() && sibling.left.isBlack() {
					sibling.color, sibling.right.color = RED, BLACK
					m.rotateLeft(sibling)
					sibling = node.getSibling()
				}
				sibling.color, node.parent.color = node.parent.color, BLACK
				if node.isLeft() {
					sibling.right.color = BLACK
					m.rotateLeft(node.parent)
				} else {
					sibling.left.color = BLACK
					m.rotateRight(node.parent)
				}
			}
		}
	}
}

// 在树中对节点进行左旋，左旋时注意：左旋节点一定要有右儿子
//
// 旋转状态如下，现在是对b进行左旋
//
//	  |                                  |
//	  b                                  d
//	 / \         rotateLeft             / \
//	a   d       ------------->         b   e
//	   / \                            / \
//	  c   e                          a   c
//
// 旋转前后的中序遍历结果是相同的
func (m *Map) rotateLeft(node *Node) {
	// 获得当前节点的右儿子以及当前节点的父节点
	rightChild := node.right
	parent := node.parent
	// 先更新当前节点，当前节点的父节点 注意，如果涉及到根节点要更新root
	if parent == nil {
		m.root = rightChild
	} else {
		if node.isLeft() {
			parent.left = rightChild
		} else {
			parent.right = rightChild
		}
	}
	rightChild.parent, node.parent = parent, rightChild
	// 更换对应子节点信息将当前节点的右儿子替换成右儿子的左儿子, 为了避免出错, 先更改对应的父节点指向再更改当前右儿子的孩子节点信息
	node.right = rightChild.left
	node.right.parent = node
	rightChild.left = node
}

// 在树中对节点进行左旋，右旋时注意：右旋节点一定要有左儿子
//
// 旋转状态如下，现在是对d进行右旋
//
//	    |                                  |
//	    d                                  b
//	   / \         rotateLeft             / \
//	  b   e       ------------->         a   d
//	 / \                            		/ \
//	a   c                          		   c   e
//
// 旋转前后的中序遍历结果是相同的
func (m *Map) rotateRight(node *Node) {
	// 获得左儿子信息以及当前节点的父节点
	leftChild := node.left
	parent := node.parent
	// 先更新当前节点的父节点的指向 注意，如果涉及到根节点要更新root
	if parent == nil {
		m.root = leftChild
	} else {
		if node.isLeft() {
			parent.left = leftChild
		} else {
			parent.right = leftChild
		}
	}
	leftChild.parent, node.parent = parent, leftChild
	// 再依次关系当前节点的的左儿子指向，左儿子反指回，前左儿子的右儿子指向
	node.left = leftChild.right
	node.left.parent = node
	leftChild.right = node
}

// 使用chan遍历map
func (m *Map) ran(node *Node, ch chan<- Pair) {
	if !node.isLeaf() {
		m.ran(node.left, ch)
		ch <- Pair{
			Key: node.key,
			Val: node.val,
		}
		m.ran(node.right, ch)
	}
}
