package main

import (
	"bufio"
	"ja_pierdole/polka"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	HOST   = "0.0.0.0"
	PORT   = "8877"
	TYPE   = "tcp"
	LEVELS = 1000
	FLAG   = "surctf_polish_notation_invented_by_bobr"
)

func main() {
	polka.Seed(time.Now().UnixNano())

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

	conn.Write([]byte("Wykreśl znaczenie wyrażenia kurwa mać!\nWyrażenia używają tylko operatorów: '+', '-', '*'.\nWszystkie liczby są całkowite!\n"))

	for i := 0; i < LEVELS; i++ {
		conn.Write([]byte("[EXP: " + strconv.Itoa(i+1) + "] "))

		exp := polka.GenerateExpr()
		conn.Write([]byte(exp.String + "\n"))

		conn.Write([]byte("Wynik: "))
		bytesBuff, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Println("Client left")
			return
		}

		result, err := strconv.Atoi(string(bytesBuff[:len(bytesBuff)-1]))
		if err != nil || result != exp.Result() {
			conn.Write([]byte("Zły wynik głupi ruski Iwan. Zacznijcie, kurwa, od nowa!\n"))
			return
		}
	}

	conn.Write([]byte("BRAWO! BRAWO! Podjąłeś właściwą decyzję, trzymaj swoją pierdoloną flagę!\n"))
	conn.Write([]byte(FLAG + "\n"))
}
