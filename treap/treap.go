// Package treap implements a treap (cartesian tree).
package treap

import (
	"cmp"
	"iter"
	"math/rand/v2"
)

// Treap is an ordered map implemented using treap.
// Its zero value is empty and ready to use.
type Treap[K cmp.Ordered, V any] struct {
	root *Node[K, V]
}

// Set inserts key with the given value.
// If key is already present, its value is overwritten.
func (t *Treap[K, V]) Set(key K, val V) {
	l, x, r := split(t.root, key)
	if x != nil {
		x.val = val
	} else {
		x = &Node[K, V]{key: key, val: val, priority: rand.Int(), sz: 1}
	}
	t.root = merge(l, merge(x, r))
}

// Delete removes key and returns its node, or nil if key was absent.
func (t *Treap[K, V]) Delete(key K) *Node[K, V] {
	l, x, r := split(t.root, key)
	t.root = merge(l, r)
	return x
}

// Get returns the node for key, or nil if key is absent.
func (t *Treap[K, V]) Get(key K) *Node[K, V] {
	return t.root.search(key)
}

// Min returns the node with the smallest key, or nil if the map is empty.
func (t *Treap[K, V]) Min() *Node[K, V] {
	return t.root.min()
}

// Max returns the node with the largest key, or nil if the map is empty.
func (t *Treap[K, V]) Max() *Node[K, V] {
	return t.root.max()
}

// Ceil returns the node with the smallest key >= key, or nil if there is none.
func (t *Treap[K, V]) Ceil(key K) *Node[K, V] {
	return t.root.ceil(key)
}

// Floor returns the node with the largest key <= key, or nil if there is none.
func (t *Treap[K, V]) Floor(key K) *Node[K, V] {
	return t.root.floor(key)
}

// Len returns the number of elements in the treap.
func (t *Treap[K, V]) Len() int {
	return t.root.size()
}

// Rank returns the number of keys < key.
func (t *Treap[K, V]) Rank(key K) int {
	return t.root.rank(key)
}

// Nth returns the node at the given rank.
func (t *Treap[K, V]) Nth(rank int) *Node[K, V] {
	return t.root.nth(rank)
}

// All returns all nodes, in ascending key order.
func (t *Treap[K, V]) All() iter.Seq[*Node[K, V]] {
	return func(yield func(*Node[K, V]) bool) {
		t.root.push(yield)
	}
}

// Range returns all nodess whose keys lie in [lo, hi], in ascending key order.
func (t *Treap[K, V]) Range(lo, hi K) iter.Seq[*Node[K, V]] {
	return func(yield func(*Node[K, V]) bool) {
		t.root.pushRange(yield, lo, hi)
	}
}

// Node is an entry in a [Treap].
type Node[K cmp.Ordered, V any] struct {
	key         K
	val         V
	priority    int
	left, right *Node[K, V]
	sz          int // Size of the subtree.
}

// Key returns the node's key.
func (n *Node[K, V]) Key() K { return n.key }

// Val returns the node's value.
func (n *Node[K, V]) Val() V { return n.val }

func (n *Node[K, V]) search(key K) *Node[K, V] {
	if n == nil {
		return nil
	}
	switch c := cmp.Compare(key, n.key); {
	case c == 0:
		return n
	case c < 0:
		return n.left.search(key)
	default: // c > 0
		return n.right.search(key)
	}
}

func (n *Node[K, V]) min() *Node[K, V] {
	if n == nil || n.left == nil {
		return n
	}
	return n.left.min()
}

func (n *Node[K, V]) max() *Node[K, V] {
	if n == nil || n.right == nil {
		return n
	}
	return n.right.max()
}

func (n *Node[K, V]) floor(key K) *Node[K, V] {
	if n == nil {
		return nil
	}
	switch c := cmp.Compare(key, n.key); {
	case c == 0:
		return n
	case c < 0:
		return n.left.floor(key)
	default: // n is a candidate, but the right subtree may hold a larger one
		if t := n.right.floor(key); t != nil {
			return t
		}
		return n
	}
}

func (n *Node[K, V]) rank(key K) int {
	if n == nil {
		return 0
	}
	switch c := cmp.Compare(key, n.key); {
	case c == 0:
		return n.left.size()
	case c < 0:
		return n.left.rank(key)
	default:
		return 1 + n.left.size() + n.right.rank(key)
	}
}

func (n *Node[K, V]) nth(rank int) *Node[K, V] {
	if n == nil {
		return nil
	}
	switch ls := n.left.size(); {
	case rank == ls:
		return n
	case rank < ls:
		return n.left.nth(rank)
	default:
		return n.right.nth(rank - ls - 1)
	}
}

func (n *Node[K, V]) ceil(key K) *Node[K, V] {
	if n == nil {
		return nil
	}
	switch c := cmp.Compare(key, n.key); {
	case c == 0:
		return n
	case c > 0:
		return n.right.ceil(key)
	default: // n is a candidate, but the left subtree may hold a smaller one
		if t := n.left.ceil(key); t != nil {
			return t
		}
		return n
	}
}

// push pushes all elements to the yield function.
func (n *Node[K, V]) push(yield func(*Node[K, V]) bool) bool {
	if n == nil {
		return true
	}
	return n.left.push(yield) && yield(n) && n.right.push(yield)
}

func (n *Node[K, V]) pushRange(yield func(*Node[K, V]) bool, lo, hi K) bool {
	if n == nil {
		return true
	}
	switch {
	case cmp.Less(n.key, lo):
		return n.right.pushRange(yield, lo, hi)
	case cmp.Less(hi, n.key):
		return n.left.pushRange(yield, lo, hi)
	default:
		return n.left.pushRange(yield, lo, hi) && yield(n) && n.right.pushRange(yield, lo, hi)
	}
}

func (n *Node[K, V]) update() {
	n.sz = 1 + n.left.size() + n.right.size()
}

func (n *Node[K, V]) size() int {
	if n == nil {
		return 0
	}
	return n.sz
}

// merge combines two treaps into one.
//
// It requires that every key in l is less than every key in r.
func merge[K cmp.Ordered, V any](l, r *Node[K, V]) *Node[K, V] {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	// The node with the higher priority becomes the root.
	if l.priority >= r.priority {
		l.right = merge(l.right, r)
		l.update()
		return l
	}
	r.left = merge(l, r.left)
	r.update()
	return r
}

// split partitions t into:
//   - l: all keys < key
//   - x: the node with key, if present
//   - r: all keys > key
func split[K cmp.Ordered, V any](t *Node[K, V], key K) (l, x, r *Node[K, V]) {
	if t == nil {
		return nil, nil, nil
	}
	switch c := cmp.Compare(t.key, key); {
	case c < 0:
		rl, x, rr := split(t.right, key)
		t.right = rl
		t.update()
		return t, x, rr
	case c > 0:
		ll, x, lr := split(t.left, key)
		t.left = lr
		t.update()
		return ll, x, t
	default:
		l, r := t.left, t.right
		t.left, t.right = nil, nil
		t.update()
		return l, t, r
	}
}
