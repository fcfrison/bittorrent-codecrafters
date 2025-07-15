package main

import (
	"errors"
	"net/http"
	"net/url"
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
	info_hash  []byte
	peer_id    []byte
	port       int
	uploaded   int
	downloaded int
	left       int
	compact    int
}

func NewBitTorrentTrackerClient() *BitTorrentTrackerClient {
	transport := &http.Transport{
		MaxConnsPerHost:       1,
		ResponseHeaderTimeout: 10 * time.Second,
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
func (c *BitTorrentTrackerClient) discoverPeers(params *DiscoverPeersParams) (string, error) {
	if params == nil {
		return "", errors.New("error: the params pointer is null")
	}
	request, err := http.NewRequest("GET", c.url.Host, nil)
	if err != nil {
		return "", err
	}
}
