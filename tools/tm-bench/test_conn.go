package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"os"
)

type line struct {
	target string
	conn   []*websocket.Conn
}

func NewLine(target string) *line {
	return &line{
		target: target,
		conn:   make([]*websocket.Conn, 1),
	}
}
func (l *line) start() error {

	for i := 0; i < 1; i++ {
		c, _, err := connect(l.target)
		if err != nil {
			return nil
		}
		l.conn[i] = c
	}
	return nil

}
func (l *line) Message(message rpctypes.RPCRequest) error {
	c := l.conn[0]
	err := c.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	l := NewLine("127.0.0.1:26657")
	if err := l.start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
