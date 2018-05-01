package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/eltrufas/rltetris"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	tetris := rltetris.CreateTetris()
	buf := make([]byte, 1024)
	sendBuf := new(bytes.Buffer)

	for !tetris.Terminal() {
		sendBuf.Reset()
		state := tetris.GetByteState()
		conn.Write(state)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connección cerrada")
				conn.Close()
				return
			}
			fmt.Println("Error leyendo", err.Error())
		}

		action := 0
		if buf[0] != 0 {
			action = 1 << (buf[0] - 1)
		}

		r := float64(tetris.Transition(uint32(action)))
		done := tetris.Terminal()

		binary.Write(sendBuf, binary.LittleEndian, r)
		binary.Write(sendBuf, binary.LittleEndian, done)

		conn.Write(sendBuf.Bytes())
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connección cerrada")
				conn.Close()
				return
			}
			fmt.Println("Error leyendo", err.Error())
		}
	}

	state := tetris.GetByteState()
	conn.Write(state)

	fmt.Println("Episodio terminado")

	conn.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}

	defer l.Close()
	fmt.Println("Escuchando")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error aceptando:", err.Error())
		}

		go handleConnection(conn)
	}
}
