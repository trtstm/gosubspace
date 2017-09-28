package main

import (
	"net"
	"strconv"
)

type ServerConnection struct {
	conn *net.UDPConn
}

func NewServerConnection(host string, port uint) (*ServerConnection, error) {
	sc := &ServerConnection{}

	serverAddr, err := net.ResolveUDPAddr("udp", host+":"+strconv.FormatUint(uint64(port), 10))
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return nil, err
	}

	sc.conn = conn

	return sc, nil
}

func (sc *ServerConnection) Stop() {
	sc.conn.Close()
}

func (sc *ServerConnection) SendBytes(buffer []byte) error {
	_, err := sc.conn.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}
