package components

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sort"
)

type KeysInt []int

func (a KeysInt) Len() int           { return len(a) }
func (a KeysInt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a KeysInt) Less(i, j int) bool { return a[i] < a[j] }
func (a *KeysInt) Push(val ...int) {
	*a = append(*a, val...)
}
func (a *KeysInt) Pop() int {
	thisArr := *a
	val := thisArr[thisArr.Len()-1]
	*a = thisArr[:thisArr.Len()-1]
	return val
}
func (a *KeysInt) Insert(index int, val ...int) {
	thisArr := *a
	if index == 0 {
		*a = append(KeysInt(val), thisArr...)
	} else if index >= thisArr.Len() {
		a.Push(val...)
	} else {
		end := append(KeysInt{}, thisArr[index:]...)
		start := append(thisArr[:index], val...)
		*a = append(start, end...)
	}
}

type KeysData [][]byte

func (a KeysData) Len() int           { return len(a) }
func (a KeysData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a KeysData) Less(i, j int) bool { return bytes.Compare(a[i], a[j]) == -1 }
func (a *KeysData) Push(val ...[]byte) {
	*a = append(*a, val...)
}
func (a *KeysData) Pop() []byte {
	thisArr := *a
	val := thisArr[thisArr.Len()-1]
	*a = thisArr[:thisArr.Len()-1]
	return val
}

func (a *KeysData) Insert(index int, val ...[]byte) {
	thisArr := *a
	if index == 0 {
		*a = append(KeysData(val), thisArr...)
	} else if index >= thisArr.Len() {
		a.Push(val...)
	} else {
		end := append(KeysData{}, thisArr[index:]...)
		start := append(thisArr[:index], val...)
		*a = append(start, end...)
	}
}

type Nodes []*TreeNode

func (a Nodes) Len() int { return len(a) }

// func (a Nodes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a Nodes) Less(i, j int) bool { return a[i].Data < a[j].Data }
func (a *Nodes) Push(val ...*TreeNode) {
	*a = append(*a, val...)
}
func (a *Nodes) Pop() *TreeNode {
	thisArr := *a
	val := thisArr[thisArr.Len()-1]
	*a = thisArr[:thisArr.Len()-1]
	return val
}
func (a *Nodes) Insert(index int, val ...*TreeNode) {
	thisArr := *a
	if index == 0 {
		*a = append(Nodes(val), thisArr...)
	} else if index >= thisArr.Len() {
		a.Push(val...)
	} else {
		end := append(Nodes{}, thisArr[index:]...)
		start := append(thisArr[:index], val...)
		*a = append(start, end...)
	}
}

func Merge[T any](s ...[]T) (slice []T) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
	default:
		s1 := s[0]
		s2 := Merge(s[1:]...) //...将数组元素打散
		slice = make([]T, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
	}
	return
}

const TreeLimit = 2

type TreeNode struct {
	treeLimit   int
	KeyNum      int         //存放KEY数量
	Key         KeysInt     //实际数据
	Parent      *TreeNode   //父节点
	ParentIndex int         //父节点位置
	Next        *TreeNode   //相连兄弟节点
	Child       Nodes       //儿子节点
	IsLeaf      bool        //是否叶子节点
	Data        interface{} //数据
}

func NewTreeNode(nodeLimit int) *TreeNode {
	if nodeLimit < TreeLimit {
		nodeLimit = TreeLimit
	}
	return &TreeNode{
		treeLimit: nodeLimit,
	}
}

// 是否满了
func (t *TreeNode) IsFull() bool {
	return t.KeyNum == t.treeLimit
}

// 插入一个键值
//
//	func (t *TreeNode) Insert(key int, data interface{}, node *TreeNode) {
//		//未满状态
//		//没有值的时候直接插入
//		if t.KeyNum == 0 {
//			if node == nil {
//				node = NewTreeNode()
//				node.Data = data
//			}
//			t.insert(key, node, 0)
//			t.ChangeParentIndex()
//			return
//		}
//
//		i := 0
//		for i < TreeLimit {
//			//没有值的时候直接插入
//			if i >= t.KeyNum {
//				//如果当前节点不是叶子节点,并且位置大于0
//				if i > 0 && !t.IsLeaf && node == nil {
//					t.Child[i-1].Insert(key, data, node)
//					return
//				}
//				if node == nil {
//					node = NewTreeNode()
//					node.Data = data
//				}
//				t.insert(key, node, i)
//				t.ChangeParentIndex()
//				return
//			}
//			//KEY值相等就把数据插入到相等KEY的数据里
//			if key == t.Key[i] {
//				if !t.IsLeaf {
//					t.Child[i].Insert(key, data, nil)
//					return
//				}
//				t.Child[i].SetValue(data)
//				return
//			}
//			//KEY值小于一个节点内的值
//			if key < t.Key[i] {
//				if !t.IsLeaf && node == nil {
//					t.Child[i].Insert(key, data, nil)
//					return
//				}
//				if node == nil {
//					node = NewTreeNode()
//					node.Data = data
//				}
//				kickKey, kickNode := t.insert(key, node, i)
//				t.ChangeParentIndex()
//				//如果要分裂新兄弟节点
//				if kickNode != nil {
//					siblingNode := t.Next
//					if siblingNode != nil {
//						t.Next.Insert(kickKey, kickNode.Data, nil)
//						return
//					}
//					t.SplitNode(kickKey, kickNode.Data, kickNode)
//					return
//				}
//
//				return
//			}
//			i++
//		}
//		//如果没有找到小于的值,就插入到最后一个子节点里面
//		if t.IsLeaf || node != nil {
//			t.SplitNode(key, data, node)
//			return
//		}
//		t.Child[i-1].Insert(key, data, node)
//		return
//	}
func (t *TreeNode) Insert(key int, data interface{}, node *TreeNode) {
	//没有值的时候直接插入
	if t.Key.Len() == 0 {
		if node == nil {
			node = NewTreeNode(t.treeLimit)
			node.Data = data
		}
		t.insert(key, node, 0)
		t.ChangeParentIndex()
		return
	}

	index := sort.Search(t.Key.Len(), func(i int) bool {
		return t.Key[i] >= key
	})

	//没有找到值,并且已满
	if index == t.Key.Len() {
		//如果没有找到小于的值,就插入到最后一个子节点里面
		if t.KeyNum == t.treeLimit && (t.IsLeaf || node != nil) {
			t.SplitNode(key, data, node)
			return
		}
		if !t.IsLeaf && node == nil {
			t.Child[index-1].Insert(key, data, node)
			return
		}

		if node == nil {
			node = NewTreeNode(t.treeLimit)
			node.Data = data
		}
		t.insert(key, node, index)
		t.ChangeParentIndex()
		return
	}

	//KEY值相等就把数据插入到相等KEY的数据里
	if t.Key[index] == key {
		if !t.IsLeaf {
			t.Child[index].Insert(key, data, node)
			return
		}
		t.Child[index].SetValue(data)
		return
	}

	//KEY值小于找到的值
	if !t.IsLeaf && node == nil {
		t.Child[index].Insert(key, data, nil)
		return
	}
	if node == nil {
		node = NewTreeNode(t.treeLimit)
		node.Data = data
	}
	kickKey, kickNode := t.insert(key, node, index)
	t.ChangeParentIndex()
	//如果要分裂新兄弟节点
	if kickNode != nil {
		if t.Next != nil {
			t.Next.Insert(kickKey, kickNode.Data, nil)
			return
		}
		t.SplitNode(kickKey, kickNode.Data, kickNode)
	}
}
func (t *TreeNode) insert(key int, node *TreeNode, index int) (outKey int, outNode *TreeNode) {
	//如果已满弹出最后的值
	node.Parent = t
	node.ParentIndex = index
	if t.KeyNum == t.treeLimit {
		outKey = t.Key.Pop()
		outNode = t.Child.Pop()
		t.KeyNum = t.Key.Len()
	}
	t.Key.Insert(index, key)
	t.Child.Insert(index, node)
	t.KeyNum = t.Key.Len()
	return
}

// 分裂兄弟节点
func (t *TreeNode) SplitNode(key int, data interface{}, node *TreeNode) {
	siblingNode := NewTreeNode(t.treeLimit)
	siblingNode.IsLeaf = t.IsLeaf
	t.Next = siblingNode
	siblingNode.Insert(key, data, node)
	if t.Parent != nil {
		t.Parent.Insert(key, nil, siblingNode)
	} else {
		parent := NewTreeNode(t.treeLimit)
		currKey, _ := t.GetRightNode()
		parent.Insert(currKey, nil, t)
		sibKey, _ := siblingNode.GetRightNode()
		parent.Insert(sibKey, nil, siblingNode)
	}
}

// 改变父节点索引值
func (t *TreeNode) ChangeParentIndex() {
	if t.Parent != nil {
		lastKey, _ := t.GetRightNode()
		t.Parent.Key[t.ParentIndex] = lastKey
		t.Parent.ChangeParentIndex()
	}
}

// 得到最右节点
func (t *TreeNode) GetRightNode() (key int, node *TreeNode) {
	key = t.Key[t.KeyNum-1]
	node = t.Child[t.KeyNum-1]
	return
}

// 得到最左节点
func (t *TreeNode) GetLeftNode() (key int, node *TreeNode) {
	key = t.Key[0]
	node = t.Child[0]
	return
}

// 设置节点数据
func (t *TreeNode) SetValue(data interface{}) {
	t.Data = fmt.Sprintf("%v,%v", t.Data, data)
}

// 搜索索引值,得到存储的数据
func (t *TreeNode) SearchIndex(idx int) (interface{}, error) {
	index := sort.Search(t.Key.Len(), func(i int) bool {
		return t.Key[i] >= idx
	})
	if index == t.Key.Len() {
		return nil, errors.New("not found index")
	}

	if !t.IsLeaf {
		return t.Child[index].SearchIndex(idx)
	}
	if t.Key[index] == idx {
		return t.Child[index].Data, nil
	} else {
		return nil, errors.New("not found index")
	}
}

func (t *TreeNode) GetBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// b+tree
type BTreePlus struct {
	rootNode *TreeNode //根节点
}

func NewBTreePlus(limit int) *BTreePlus {
	root := NewTreeNode(limit)
	root.IsLeaf = true
	return &BTreePlus{
		rootNode: root,
	}
}

func (b *BTreePlus) Insert(key int, data interface{}) {
	b.rootNode.Insert(key, data, nil)
	if b.rootNode.Parent != nil {
		b.rootNode = b.rootNode.Parent
	}
}

func (b *BTreePlus) Remove(key int) {

}

func (b *BTreePlus) Search(key int) (interface{}, error) {
	return b.rootNode.SearchIndex(key)
}

func (b *BTreePlus) Print() {
	tmp := b.rootNode
	for !tmp.IsLeaf {
		next := tmp
		for next != nil {
			fmt.Print(next.Key)
			next = next.Next
		}
		fmt.Println("")
		tmp = tmp.Child[0]
	}
	next := tmp
	for next != nil {
		fmt.Print("[")
		for i, v := range next.Key {
			fmt.Printf("%v(%v)", v, next.Child[i].Data)
		}
		fmt.Print("]", " ")
		next = next.Next
	}
	fmt.Println("")
}
