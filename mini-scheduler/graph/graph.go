package graph

type Graph struct {
	node  map[string]struct{}
	edges map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		node:  make(map[string]struct{}),
		edges: make(map[string][]string),
	}
}

func (g *Graph) AddTask(node string) {
	g.node[node] = struct{}{}
}
func (g *Graph) Adddependency(dependency string, task string) {
	g.edges[dependency] = append(g.edges[dependency], task)
}
func (g *Graph) Parents(node string) []string {
	var parents []string
	for parent, children := range g.edges {
		for _, child := range children {
			if child == node {
				parents = append(parents, parent)
			}
		}
	}
	return parents
}
func (g *Graph) Children(node string) []string {
	var children []string
	if val, ok := g.edges[node]; ok {
		children = val
	}
	return children
}
