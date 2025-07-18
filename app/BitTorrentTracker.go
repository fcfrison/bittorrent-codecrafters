package main

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BitTorrentTrackerClient struct {
	client *http.Client
	url    *url.URL
}
type BitTorrentTrackerApi interface {
	GetPeers()
}

type DiscoverPeersParams struct {
	info_hash  [20]byte
	peer_id    [20]byte
	uploaded   int
	downloaded int
	left       int
	compact    int
}

func NewDiscoverPeersParamsStruct(info_hash [20]byte, peer_id [20]byte,
	uploaded int, downloaded int, left int, compact int) *DiscoverPeersParams {
	return &DiscoverPeersParams{
		info_hash:  info_hash,
		peer_id:    peer_id,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       left,
		compact:    compact,
	}
}

var port int = 6881

func NewBitTorrentTrackerClient() *BitTorrentTrackerClient {
	localAddress := &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}
	dialer := &net.Dialer{
		LocalAddr: localAddress,
	}
	transport := &http.Transport{
		MaxConnsPerHost:       1,
		ResponseHeaderTimeout: 10 * time.Second,
		DialContext:           dialer.DialContext,
	}
	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}
	return &BitTorrentTrackerClient{
		client: client,
	}
}

func (c *BitTorrentTrackerClient) SetUrl(baseUrl string) error {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	c.url = parsedUrl
	return nil
}
func (c *BitTorrentTrackerClient) DiscoverPeers(params *DiscoverPeersParams) ([]byte, error) {
	if params == nil {
		return nil, errors.New("error: the params pointer is null")
	}
	if c.url == nil {
		return nil, errors.New("error: no url was setted yet")
	}

	request, err := http.NewRequest("GET", "https://"+c.url.Host+c.url.Path, nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("info_hash", string(params.info_hash[:]))
	q.Add("peer_id", string(params.peer_id[:]))
	q.Add("port", strconv.Itoa(port))
	q.Add("uploaded", strconv.Itoa(params.uploaded))
	q.Add("downloaded", strconv.Itoa(params.downloaded))
	q.Add("left", strconv.Itoa(params.left))
	q.Add("compact", strconv.Itoa(params.compact))
	request.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, errors.New("error: the tracker returned a status code that is outside the range [200,300)")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
