package treekeys

// Still need:
// * DHKeyGen
// * KeyExchangeKeyGen
// * KeyExchange
// * MAC

// Tree Management functions

type PrivateKey [32]byte
type GroupElement [32]byte

type Tree interface {
	Size() int
	PublicKey(x, y int) GroupElement
	Subtree(x, y int) Tree
}

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Value PrivateKey
}

func (t TreeNode) IsLeaf() bool {
	return t.Left == nil && t.Right == nil
}

func (t TreeNode) Size() int {
	if t.IsLeaf() {
		return 1
	}

	return t.Left.Size() + t.Right.Size()
}

func CreateTree(λ []PrivateKey) *TreeNode {
	n := len(λ)
	if n == 1 {
		return &TreeNode{Left: nil, Right: nil, Value: λ[0]}
	}

	m := pow2(n)
	L := CreateTree(λ[0:m])
	R := CreateTree(λ[m:n])

	k := ι(Exp(PK(L.Value), R.Value))
	return &TreeNode{Left: L, Right: R, Value: k}
}

func Copath(T *TreeNode, i int) []GroupElement {
	if T.IsLeaf() {
		return []GroupElement{}
	}

	m := pow2(T.Size())

	var key GroupElement
	var remainder []GroupElement
	if i < m {
		key = PK(T.Right.Value)
		remainder = Copath(T.Left, i)
	} else {
		key = PK(T.Left.Value)
		remainder = Copath(T.Right, i-m)
	}

	return append([]GroupElement{key}, remainder...)
}

func PathNodeKeys(λ PrivateKey, P []GroupElement) []PrivateKey {
	nks := make([]PrivateKey, len(P)+1)
	nks[len(P)] = λ
	for n := len(P) - 1; n >= 0; n -= 1 {
		nks[n] = ι(Exp(P[n], nks[n+1]))
	}
	return nks
}