package hashring

// Node define the interface required for a Node in the Tree
type Node interface {
	Key() uint64
	Value() Node
	Left() Node
	Right() Node
}

// Tree define the interface required for the Tree
type Tree interface {
	Root() Node
	Put(key uint64, n Node)
	Remove(key uint64)
	Get(key uint64) Node
}

func search(key uint64, n Node) Node {
	if n == nil {
		return nil
	}

	left := n.Left()
	right := n.Right()
	k := n.Key()
	switch {
	case key == k:
		return n
	case key < k && left == nil:
		return n
	case key > k && right == nil:
		return n
	case key < k && left != nil:
		return search(key, left)
	case key > k && right != nil:
		return search(key, right)
	}
	return nil
}
