package bottle

import (
	"strings"
)

type node struct {
	pattern string
	part string
	children []*node
	isWild bool // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) isEnd() bool {
	return len(n.children) == 0
}

func (n *node) firstMatch(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// All matched children nodes
func (n *node) allMatch(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.firstMatch(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height + 1)
}

func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.allMatch(part)
	for _, child := range children {
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	}

	return nil
}
