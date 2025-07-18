package main

import (
	"errors"
	"fmt"
)

type Peer struct {
	ip   []byte
	port int
}

func (p Peer) StrRepr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", p.ip[0], p.ip[1], p.ip[2], p.ip[3], p.port)
}
func parsePeers(peers []byte) ([]Peer, error) {
	if len(peers) == 0 {
		return nil, errors.New("error: no peer was indicated")
	}
	var ip []byte
	var port []byte
	var ipIsComplete, portIsComplete bool
	peerSlice := make([]Peer, 0)
	for i, val := range peers {
		j := i % 6
		if j == 0 {
			ip = make([]byte, 4)
			port = make([]byte, 2)
			ipIsComplete = false
			portIsComplete = false
		}
		if j <= 3 {
			ip[j] = val
			ipIsComplete = j == 3
		} else {
			port[j%4] = val
			portIsComplete = j%4 == 1
		}
		if ipIsComplete && portIsComplete {
			portInt := int(port[0])<<8 | int(port[1])
			newPeer := Peer{
				ip:   ip,
				port: portInt,
			}
			peerSlice = append(peerSlice, newPeer)
		}
	}
	return peerSlice, nil
}
