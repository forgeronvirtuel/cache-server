package cache

import (
	"fmt"
	"sort"
	"strings"
)

func splitAndTrim(path string) []string {
	keys := strings.Split(path, "/")
	if keys[0] == "" {
		keys = keys[1:]
	}
	if keys[len(keys)-1] == "" {
		keys = keys[:len(keys)-1]
	}
	return keys
}

type Root struct {
	childs []*Node
}

func NewRoot() *Root {
	return &Root{childs: []*Node{}}
}

func (r *Root) Root(path string) *Root {
	keys := splitAndTrim(path)
	if r == nil {
		r = &Root{
			childs: []*Node{},
		}
	}
	for _, n := range r.childs {
		if n.name == keys[0] {
			n.append(keys[1:])
			return r
		}
	}
	n := &Node{name: keys[0]}
	n.append(keys[1:])
	r.childs = append(r.childs, n)
	return r
}

func (r *Root) GetAllPaths(path string) []string {
	if path == "" {
		var paths []string
		for _, c := range r.childs {
			paths = append(paths, c.getAllPaths("")...)
		}
		return paths
	}

	keys := splitAndTrim(path)
	for _, c := range r.childs {
		if c.name == keys[0] {
			n := c.searchForNode(keys[1:])
			return n.getAllPaths("")
		}
	}
	return nil
}

type Node struct {
	name   string
	childs []*Node
}

//func (n *Node) Root(path string) *Node {
//	keys := splitAndTrim(path)
//	if n == nil {
//		n = &Node{
//			name:   keys[0],
//			childs: nil,
//		}
//		n.append(keys[1:])
//		return n
//	} else {
//		n.append(keys)
//	}
//	return n
//}

func (n *Node) String() string {
	return n.name
}

func (n *Node) append(keys []string) {
	// Stop condition
	if len(keys) == 0 {
		return
	}

	// Try to find a path for the given key
	for _, v := range n.childs {
		if v.name == keys[0] {
			v.append(keys[1:])
			return
		}
	}

	// If no path has been found, create it
	node := &Node{
		name:   keys[0],
		childs: nil,
	}
	n.childs = append(n.childs, node)
	sort.Slice(n.childs, func(i, j int) bool {
		return n.childs[i].name < n.childs[j].name
	})
	node.append(keys[1:])
}

func (n *Node) searchForNode(keys []string) *Node {
	if len(keys) == 0 {
		return n
	}
	for _, c := range n.childs {
		if c.name == keys[0] {
			return n.searchForNode(keys[1:])
		}
	}
	return nil
}

func (n *Node) getAllPaths(path string) []string {
	path = fmt.Sprintf("%s/%s", path, n.name)
	if len(n.childs) == 0 {
		return []string{path}
	}

	var allpaths []string
	for _, c := range n.childs {
		paths := c.getAllPaths(path)
		allpaths = append(allpaths, paths...)
	}
	return allpaths
}
