# -*- coding: utf-8 -*-
from pwn import *

PATH = -5
START = -4
FINISH = -3
WALL = -2
CELL = -1

UP = 1
DOWN = 2
RIGHT = 3
LEFT = 4

class Point:
	def __init__(self, x, y):
		self.x, self.y = x, y


# Парсим лабиринт в двумерный массив, эмодзи заменяем на константы: CELL, WALL, FINISH, START
def parse_maze(r):
	maze = []
	for i in range(31):
		line = r.recvline().decode("utf-8").replace("  ", " ").replace("\n", "")

		binary = []
		for c in line:
			if c == " ":
				binary.append(CELL)
			elif c == "🟦":
				binary.append(WALL)
			elif c == "⛳":
				binary.append(FINISH)
			elif c == "😎":
				binary.append(START)
		maze.append(binary)

	return maze

# Реализует шаг распространения волны из точки src, устанавливая соседям значение val
def init_wave(trace, src, val):
	changed = False
	if src.x < len(trace[0]) - 1:
		if trace[src.y][src.x + 1] == CELL:
			trace[src.y][src.x + 1] = val
			changed = True

	if src.x > 0:
		if trace[src.y][src.x - 1] == CELL:
			trace[src.y][src.x - 1] = val
			changed = True

	if src.y > 0:
		if trace[src.y - 1][src.x] == CELL:
			trace[src.y - 1][src.x] = val
			changed = True

	if src.y < len(trace) - 1:
		if trace[src.y + 1][src.x] == CELL:
			trace[src.y + 1][src.x] = val
			changed = True

	return changed

# Распространяет волну на один шаг для каждой точки со значением равным step
def spread_wave(trace, step):
	changed = False
	for y in range(len(trace)):
		for x in range(len(trace[y])):
			if trace[y][x] == step:
				if init_wave(trace, Point(x, y), step+1):
					changed = True

	return changed

# Выбирает соседа клетки с минимальным значением.
# Возвращает позицию соседа и движение чтобы перейти к нему.
# (Для оси Y движения инвертированы, так-как потом значения используются для составления пути от src к dst)
def get_min_neighbour(trace, cell):
	min_n = cell
	min_val = trace[cell.y][cell.x]
	move = None

	if cell.x > 0 and trace[cell.y][cell.x - 1] not in [CELL, WALL, PATH]:
		if trace[cell.y][cell.x - 1] < min_val:
			min_n = Point(cell.x - 1, cell.y)
			min_val = trace[min_n.y][min_n.x]
			move = RIGHT

	if cell.x < len(trace[0]) - 1 and trace[cell.y][cell.x + 1] not in [CELL, WALL, PATH]:
		if trace[cell.y][cell.x + 1] < min_val:
			min_n = Point(cell.x + 1, cell.y)
			min_val = trace[min_n.y][min_n.x]
			move = LEFT

	if cell.y > 0 and trace[cell.y - 1][cell.x] not in [CELL, WALL, PATH]:
		if trace[cell.y - 1][cell.x] < min_val:
			min_n = Point(cell.x, cell.y - 1)
			min_val = trace[min_n.y][min_n.x]
			move = DOWN

	if cell.y < len(trace) - 1 and trace[cell.y + 1][cell.x] not in [CELL, WALL, PATH]:
		if trace[cell.y + 1][cell.x] < min_val:
			min_n = Point(cell.x, cell.y + 1)
			min_val = trace[min_n.y][min_n.x]
			move = UP

	return min_n, move

# Находит обратный путь от dst в src. В командах UP, DOWN, LEFT, RIGHT
# Возвращает развернутый путь, т.е. путь от src до dst.
def backtrace_path(trace, src, dst):
	moves = []
	c = 0
	while dst.x != src.x or dst.y != src.y:
		min_n, move = get_min_neighbour(trace, dst)
		moves.append(move)
		trace[dst.y][dst.x] = PATH
		dst = min_n

	return moves[::-1]

# Реализует алгоритм волновой трассировки. 
# Инициирует волну из точки src и распространяет до тех пор пока не достигнет dst,
# либо не покроет все досигаемые клетки.
# В конце формиру
def trace_maze(maze):
	dst, src = None, None
	for y, line in enumerate(maze):
		for x, c in enumerate(line):
			if c == START:
				src = Point(x, y)
			
			if c == FINISH:
				dst = Point(x, y)

		if src and dst:
			break

	trace = [row[:] for row in maze]
	trace[dst.y][dst.x] = CELL

	init_wave(trace, src, 1)
	step = 1
	
	while spread_wave(trace, step):
		step += 1
		if trace[dst.y][dst.x] != CELL:
			break


	path = backtrace_path(trace, src, dst)
	return trace, path

# Выводит лабиринт с трассировкой.
def print_traced_maze(maze, trace):
	for y, line in enumerate(maze):
		for x, c in enumerate(line):
			if c == WALL:
				print("🟦", end="")
			elif c == CELL:
				if trace[y][x] != CELL:
					if trace[y][x] == PATH:
						print("🟩", end="")
					else:
						print("%02d" % trace[y][x], end="")
				else:
					print("  ", end="")
			elif c == START:
				print("😎", end="")
			else:
				print("⛳️", end="")

		print()


r = remote("185.104.115.19", 9966)
# r = remote("localhost", 9966)
r.recvuntil("200 команд.\n".encode("utf-8"))

c = 1
try:
	while True:
		print("MAZE: ", c)
		r.recvuntil(b"]\n")
		maze = parse_maze(r)
		trace, path = trace_maze(maze)
		print_traced_maze(maze, trace)
		
		r.sendline(''.join(map(str, path)).encode("utf-8"))


		# Пропускаем вывод пройденного пути сервером
		for i in range(31):
			r.recvline()

		c += 1
except:
	r.interactive()
finally:
	r.close()