package main

import (
	"bufio"
	"context"
	"net"
	"strconv"
	"time"
)

type ClientConfig struct {
	Address       string
	Dialer        *net.Dialer
	RecChanSize   int
	WriteChanSize int
}
type TcpClient interface {
	Connect() error
	Send()
	Receive()
}
type BitTorrentTcpClient struct {
	config  *ClientConfig
	conn    *net.Conn
	readCh  chan *ReceivedData
	writeCh chan []byte
}

func NewClientConfig(urlOrIp string, port int) *ClientConfig {
	dialer := &net.Dialer{
		Timeout:       5 * time.Second,
		KeepAlive:     15 * time.Second,
		FallbackDelay: 300 * time.Millisecond,
	}
	return &ClientConfig{
		Address: urlOrIp + ":" + strconv.Itoa(port),
		Dialer:  dialer,
	}

}
func NewBitTorrentTcpClient(config *ClientConfig) (*BitTorrentTcpClient, error) {

	return &BitTorrentTcpClient{
		config: config,
	}, nil
}
func (c *BitTorrentTcpClient) Connect() error {
	conn, err := c.config.Dialer.DialContext(context.Background(), "tcp", c.config.Address)
	if err != nil {
		return err
	}
	c.conn = &conn
	c.readCh = make(chan *ReceivedData, c.config.RecChanSize)
	c.writeCh = make(chan []byte, c.config.WriteChanSize)
	return nil
}

type ReceivedData struct {
	data []byte
	size int
	err  error
}

func (c *BitTorrentTcpClient) Receive() {
	reader := bufio.NewReader(*c.conn)
	for {
		buff := make([]byte, 1024)
		size, err := reader.Read(buff)
		recData := &ReceivedData{
			data: buff,
			size: size,
			err:  err}
		c.readCh <- recData
		if err != nil {
			return
		}
	}
}
func (c *BitTorrentTcpClient) Send() {
	writter := bufio.NewWriter(*c.conn)
	for {
		buff, ok := <-c.writeCh
		if !ok {
			return
		}
		buffSize := len(buff)
		for buffSize > 0 {
			n, err := writter.Write(buff)
			if err != nil {
				return
			}
			buffSize = buffSize - n
		}
		err := writter.Flush()
		if err != nil {
			return
		}
	}
}
