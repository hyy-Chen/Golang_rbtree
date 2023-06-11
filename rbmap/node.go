package rbmap

type keyItem interface{}

type valItem interface{}

// Node 节点结构体，实现的方法都是不安全的，未进行越界判断的
type Node struct {
	key                 keyItem // 键值
	val                 valItem // 价值
	left, right, parent *Node   // 左，右指针和指向父节点的指针
	color               bool    // 节点颜色
}

const (
	// RED 		红色为1
	RED = true
	// BLACK 	黑色为0
	BLACK = false
)

// 创建叶子节点
func newLeaf() *Node {
	return &Node{
		left:   nil,
		right:  nil,
		parent: nil,
		color:  BLACK,
	}
}

// NewNode 创建红色节点
func newNode(key keyItem, val valItem) *Node {
	return &Node{
		key:    key,
		val:    val,
		left:   nil,
		right:  nil,
		parent: nil,
		color:  RED,
	}
}

// 判断是否是黑色叶子节点
func (n *Node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

// 判断是不是根节点
func (n *Node) isRoot() bool {
	return n.parent == nil
}

// 判断是不是红色节点
func (n *Node) isRed() bool {
	return n.color
}

// 判断是不是黑色节点
func (n *Node) isBlack() bool {
	return !n.color
}

// 此节点是左儿子
func (n *Node) isLeft() bool {
	return n == n.parent.left
}

// 此节点是右儿子
func (n *Node) isRight() bool {
	return n == n.parent.right
}

// 获得祖父节点
func (n *Node) getGrandParent() *Node {
	return n.parent.parent
}

// 获得兄弟节点
func (n *Node) getSibling() *Node {
	if n.isLeft() {
		return n.parent.right
	}
	return n.parent.left
}

// 获得叔父节点
func (n *Node) getUncle() *Node {
	return n.parent.getSibling()
}
