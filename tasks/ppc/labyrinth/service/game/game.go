package game

import (
	"labyrinth/maze"
	"math/rand"
)

const (
	CRASH  = "🟥"
	PATH   = "🟩"
	WALL   = "🟦"
	CELL   = "  "
	PLAYER = "😎"
	TARGET = "⛳️"
	UP     = 1
	DOWN   = 2
	RIGHT  = 3
	LEFT   = 4
)

type Point struct {
	X, Y int
}

type Game struct {
	Player Point
	Target Point
	Maze   maze.Maze
}

func Seed(n int64) {
	maze.Seed(n)
}

func (g *Game) IsWon() bool {
	return g.Player.X == g.Target.X && g.Player.Y == g.Target.Y
}

func (g *Game) ToString() string {
	s := ""
	for y, line := range g.Maze.Binary {
		for x, c := range line {
			if x == g.Player.X && y == g.Player.Y {
				s = s + PLAYER
				continue
			}

			if x == g.Target.X && y == g.Target.Y {
				s = s + TARGET
				continue
			}

			switch c {
			case maze.CELL:
				s = s + CELL
			case maze.WALL:
				s = s + WALL
			case 3:
				s = s + PATH
			case 4:
				s = s + CRASH
			}
		}
		s = s + "\n"
	}

	return s
}

func (g *Game) canMove(move int) bool {
	switch move {
	case UP:
		return g.Maze.Binary[g.Player.Y-1][g.Player.X] == 0
	case DOWN:
		return g.Maze.Binary[g.Player.Y+1][g.Player.X] == 0
	case RIGHT:
		return g.Maze.Binary[g.Player.Y][g.Player.X+1] == 0
	case LEFT:
		return g.Maze.Binary[g.Player.Y][g.Player.X-1] == 0
	}

	return false
}

func (g *Game) MovePlayer(move int) bool {
	switch {
	case move == UP && g.canMove(move):
		g.Maze.Binary[g.Player.Y][g.Player.X] = 3 // PATH BLOCK
		g.Player.Y--
	case move == DOWN && g.canMove(move):
		g.Maze.Binary[g.Player.Y][g.Player.X] = 3 // PATH BLOCK
		g.Player.Y++
	case move == RIGHT && g.canMove(move):
		g.Maze.Binary[g.Player.Y][g.Player.X] = 3 // PATH BLOCK
		g.Player.X++
	case move == LEFT && g.canMove(move):
		g.Maze.Binary[g.Player.Y][g.Player.X] = 3 // PATH BLOCK
		g.Player.X--
	default:
		switch move {
		case UP:
			g.Maze.Binary[g.Player.Y-1][g.Player.X] = 4 // CRASH BLOCK
		case DOWN:
			g.Maze.Binary[g.Player.Y+1][g.Player.X] = 4 // CRASH BLOCK
		case LEFT:
			g.Maze.Binary[g.Player.Y][g.Player.X-1] = 4 // CRASH BLOCK
		case RIGHT:
			g.Maze.Binary[g.Player.Y][g.Player.X+1] = 4 // CRASH BLOCK
		}
		return false
	}

	return true
}

func NewGame(rows, cols int) Game {
	game := Game{Player: Point{1, 1}, Maze: maze.GenerateMaze(rows, cols)}

	limitX, limitY := len(game.Maze.Binary[0]), len(game.Maze.Binary)
	var target Point
	for {
		target.X, target.Y = rand.Int()%limitX, rand.Int()%limitY

		if game.Maze.Binary[target.Y][target.X] != 0 {
			continue
		}
		if target.X == game.Player.X && target.Y == game.Player.Y {
			continue
		}
		break
	}

	game.Target = target
	return game
}
