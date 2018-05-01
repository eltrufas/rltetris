package rltetris

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/eltrufas/tetriscore"
)

type RemotePlayer struct {
	conn net.Conn
}

func CreateRemotePlayer(address string) *RemotePlayer {
	var player RemotePlayer

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("No se pudo conectar al jugador remoto", err.Error())
		return nil
	}

	player.conn = conn

	return &player
}

func (ps *RemotePlayer) GetAction(t *tetriscore.Tetris) tetriscore.InputState {
	byteState := GetByteState(t)

	ps.conn.Write(byteState)

	buf := make([]byte, 32)
	_, err := ps.conn.Read(buf)
	if err != nil {
		fmt.Println("Error leyendo acción", err.Error())
		return 0
	}

	action := 0
	if buf[0] != 0 {
		action = 1 << (buf[0] - 1)
	}

	return tetriscore.InputState(action)
}

func handleConnection(conn net.Conn) {
	tetris := CreateTetris()
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

func playListen() (chan tetriscore.Tetris, chan tetriscore.InputState) {
	l, err := net.Listen("tcp", "localhost:4040")
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

func replayListen() {
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
