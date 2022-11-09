package maze

// Implementation of Kruskal's maze generation algorithm

import (
	"labyrinth/tree"
	"math/rand"
)

type Maze struct {
	Grid       [][]tree.Node
	Binary     [][]byte
	Rows, Cols int
}

func (m *Maze) ToString(cell, wall string) string {
	s := ""
	for _, line := range m.Binary {
		for _, c := range line {
			switch c {
			case 0:
				s = s + cell
			case 1:
				s = s + wall
			}
		}
		s = s + "\n"
	}

	return s
}

func Seed(n int64) {
	rand.Seed(n)
}

const (
	WALL = 1
	CELL = 0
)

func binaryPresentation(grid [][]tree.Node, deletedEdges []tree.Edge) [][]byte {
	height, width := len(grid)*2+1, len(grid[0])*2+1

	output := make([][]byte, height)
	for y := 0; y < height; y++ {
		output[y] = make([]byte, width)

		if y%2 == 0 {
			for i := 0; i < width; i++ {
				output[y][i] = 1
			}
			continue
		}

		for x := 0; x < width; x++ {
			if x%2 == 0 {
				output[y][x] = WALL
			} else {
				output[y][x] = CELL
			}
		}
	}

	for _, edge := range deletedEdges {
		x, y := (edge.X+1)*2-1, (edge.Y+1)*2-1
		if edge.IsVertical {
			x += edge.D
		} else {
			y += edge.D
		}

		output[x][y] = CELL
	}

	return output
}

func GenerateMaze(n, m int) Maze {
	var edges []tree.Edge
	grid := make([][]tree.Node, n, n)
	id := 0
	for y := 0; y < n; y++ {
		grid[y] = make([]tree.Node, 0, m)
		for x := 0; x < m; x++ {
			grid[y] = append(grid[y], tree.Node{ID: id, X: x, Y: y})

			if x > 0 {
				edges = append(edges, tree.Edge{x, y, -1, true})
			}

			if x < m-1 {
				edges = append(edges, tree.Edge{x, y, +1, true})
			}

			if y > 0 {
				edges = append(edges, tree.Edge{x, y, -1, false})
			}

			if y < n-1 {
				edges = append(edges, tree.Edge{x, y, +1, false})
			}

			id++
		}
	}

	rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })

	var deletedEdges []tree.Edge
	// Randomly connecting cells by deleting edges
	for len(edges) > 0 {
		edge := edges[0]
		edges = edges[1:]

		x0, y0 := edge.X, edge.Y
		x1, y1 := edge.X, edge.Y
		if edge.IsVertical {
			x1 += edge.D
		} else {
			y1 += edge.D
		}

		if grid[y0][x0].IsConnectedWith(&grid[y1][x1]) {
			continue
		}

		deletedEdges = append(deletedEdges, edge)
		grid[y0][x0].Connected = append(grid[y0][x0].Connected, &grid[y1][x1])
		grid[y1][x1].Connected = append(grid[y1][x1].Connected, &grid[y0][x0])
	}

	return Maze{Grid: grid, Binary: binaryPresentation(grid, deletedEdges), Rows: n, Cols: m}
}
