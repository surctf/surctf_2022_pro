package tree

type Tree struct {
	Nodes []*Node
}

type Node struct {
	ID, X, Y   int
	Connected  []*Node
	Discovered bool
}

func dfs(root *Node, wantedNode *Node) bool {
	var visited []*Node
	visited = append(visited, root)
	discovered := make(map[int]bool)

	for len(visited) > 0 {
		n := visited[0]
		discovered[n.ID] = true

		visited = visited[1:]
		for _, connectedNode := range n.Connected {
			if _, ok := discovered[connectedNode.ID]; ok {
				continue
			}

			if connectedNode.ID == wantedNode.ID {
				return true
			}
			visited = append(visited, connectedNode)
		}
	}

	return false
}

func (n *Node) IsConnectedWith(otherNode *Node) bool {
	isConnected := dfs(n, otherNode)
	return isConnected
}

type Edge struct {
	X, Y, D    int
	IsVertical bool
}
