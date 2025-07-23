package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type ClientConfig struct {
	Address       string
	Dialer        *net.Dialer
	RecChanSize   int
	WriteChanSize int
}
type TCPClient interface {
	Connect() error
	Disconnect() error
	Send([]byte) error
	Receive() ([]byte, error)
}
type BitTorrentTcpClient struct {
	config  *ClientConfig
	conn    *net.Conn
	readCh  chan *ReceivedData
	writeCh chan []byte
	writter *bufio.Writer
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
	c.writter = bufio.NewWriter(conn)
	return nil
}

type ReceivedData struct {
	data []byte
	size int
}

func (c *BitTorrentTcpClient) Read() {
	reader := bufio.NewReader(*c.conn)
	for {
		buff := make([]byte, 1024)
		size, err := reader.Read(buff)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if size > 0 {
			c.readCh <- &ReceivedData{data: buff, size: size}
		}

	}
}
