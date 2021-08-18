package termites

type NodeRegistry interface {
	register(n *node)
	Iterate() []Node
}

type nodeRegistry struct {
	registeredNodes map[NodeId]*node
}

func newNodeRegistry() *nodeRegistry {
	return &nodeRegistry{
		registeredNodes: make(map[NodeId]*node),
	}
}

func (r *nodeRegistry) Iterate() []Node {
	var nodes []Node
	for _, n := range r.registeredNodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (r *nodeRegistry) register(n *node) {
	_, has := r.registeredNodes[n.id]
	if !has {
		r.registeredNodes[n.id] = n
	}
}
