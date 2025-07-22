package main

import (
	"context"
	"net"
	"strconv"
	"time"
)

type BitTorrentTcpClient struct {
	client  *net.Conn
	address string
}

func NewBitTorrentTcpClient(address string, port int) (*BitTorrentTcpClient, error) {
	dialer := &net.Dialer{
		Timeout:       5 * time.Second,
		KeepAlive:     15 * time.Second,
		FallbackDelay: 300 * time.Millisecond,
	}
	completeAddress := address + ":" + strconv.Itoa(port)
	conn, err := dialer.DialContext(context.Background(), "tcp", completeAddress)
	if err != nil {
		return nil, err
	}
	return &BitTorrentTcpClient{
		client:  &conn,
		address: completeAddress,
	}, nil
}
