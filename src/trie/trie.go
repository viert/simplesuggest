package trie

type Trie struct {
	root     *Node
	trieSize int
}

func NewTrie() *Trie {
	t := new(Trie)
	t.root = newNode()
	t.trieSize = EMPTY_NODE_SIZE
	return t
}

func (t *Trie) Size() int {
	return t.trieSize
}

func (t *Trie) Insert(key string, objectId string) {
	t.trieSize = t.trieSize + t.root.insert(key, objectId, 0)
}

func (t *Trie) GetAll() []string {
	return t.root.getAll()
}

func (t *Trie) GetAllWithPrefix(prefix string) []string {
	return t.root.getAllWithPrefix(prefix, 0)
}
