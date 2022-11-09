package main

/* TODO:
1. Info about controls, labyrinth always static size
2. Docker
*/

import (
	"bufio"
	"labyrinth/game"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	HOST         = "0.0.0.0"
	PORT         = "9966"
	TYPE         = "tcp"
	LEVELS       = 250
	MAX_PATH_LEN = 200
	FLAG         = "surctf_easy_m4ze_n0t_3asy_alg"
)

func main() {
	game.Seed(time.Now().UnixNano())

	log.Printf("Starting TCP server on %s:%s", HOST, PORT)
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println("Error connecting:", err.Error())
			return
		}

		log.Printf("Client %s connected\n", c.RemoteAddr().String())

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	welcomeMSG := ("😎 - твоя позиция\n⛳️ - точка в которую нужно попасть\n" +
		"Все лабиринты одного размера и удовлетворяют условию 'идеального лабиринта' a.k.a 'perfect maze'.\n" +
		"Для передвижения в лабиринте отправь мне строку из цифр: 1(Вверх), 2(Вниз), 3(Вправо), 4(Влево).\n" +
		"Например, строка '223314' будет означать 'вниз-вниз-вправо-вправо-вверх-влево'.\nМаксимальная длина пути - 200 команд.\n")
	conn.Write([]byte("[INFO]\n" + welcomeMSG))

	for i := 0; i < LEVELS; i++ {
		g := game.NewGame(15, 15)
		for {
			conn.Write([]byte("[MAZE: " + strconv.Itoa(i+1) + "]\n"))
			conn.Write([]byte(g.ToString()))

			conn.Write([]byte("Path: "))
			bytesBuff, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				log.Println("Client left")
				return
			}

			path := bytesBuff[:len(bytesBuff)-1]
			if len(path) > MAX_PATH_LEN {
				conn.Write([]byte("Слишком длинный путь! Максимальное длина - 200 команд.\n"))
				return
			}

			for _, move := range path {
				if !g.MovePlayer(int(move) - 48) {
					conn.Write([]byte(g.ToString()))
					conn.Write([]byte("Туда нельзя! Пока.\n"))
					return
				}
			}

			conn.Write([]byte(g.ToString()))

			if !g.IsWon() {
				conn.Write([]byte("Путь не привел тебя к цели((( Пока.\n"))
				return
			} else {
				break
			}
		}
		conn.Write([]byte("GOOD! RESHAI DALSHE!\n"))
	}

	conn.Write([]byte("MOLODEC! DERJY FLAG!\n"))
	conn.Write([]byte(FLAG + "\n"))
}
