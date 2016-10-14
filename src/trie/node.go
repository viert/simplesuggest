package trie

import (
	"strings"
)

const (
	EMPTY_NODE_SIZE = 8
)

type Node struct {
	objectId string
	children map[string]*Node
}

func newNode() *Node {
	n := new(Node)
	n.children = make(map[string]*Node)
	return n
}

func (n *Node) insert(key string, objectId string, stringPos int) int {

	bytesInserted := 0

	if stringPos < 0 || stringPos > len(key)-1 {
		return bytesInserted
	}

	currentLetter := string(key[stringPos])
	if _, found := n.children[currentLetter]; !found {
		n.children[currentLetter] = newNode()
		bytesInserted = bytesInserted + EMPTY_NODE_SIZE + 1 // one more key (1), on more node (EMPTY_NODE_SIZE)
	}
	if stringPos == len(key)-1 {
		if n.children[currentLetter].objectId != objectId {
			oldObjectId := n.children[currentLetter].objectId
			n.children[currentLetter].objectId = objectId
			bytesInserted = bytesInserted + len(objectId) - len(oldObjectId)
		}
	} else {
		bytesInserted = bytesInserted + n.children[currentLetter].insert(key, objectId, stringPos+1)
	}
	return bytesInserted
}

func (n *Node) getAll() []string {
	var result []string
	if n.objectId == "" {
		result = make([]string, 0)
	} else {
		result = make([]string, 1)
		result[0] = n.objectId
	}
	for _, node := range n.children {
		result = append(result, node.getAll()...)
	}
	return result
}

func (n *Node) getAllWithPrefix(prefix string, stringPos int) []string {
	if stringPos < len(prefix) {
		currentLetter := string(prefix[stringPos])

		testLetters := make([]string, 1)
		testLetters[0] = currentLetter

		additional := strings.ToUpper(currentLetter)
		if additional != currentLetter {
			testLetters = append(testLetters, additional)
		} else {
			additional = strings.ToLower(currentLetter)
			if additional != currentLetter {
				testLetters = append(testLetters, additional)
			}
		}

		results := make([]string, 0)

		for _, letter := range testLetters {
			if _, found := n.children[letter]; found {
				results = append(results, n.children[letter].getAllWithPrefix(prefix, stringPos+1)...)
			}
		}
		return results

	} else {
		return n.getAll()
	}
}
