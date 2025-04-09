package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"namereg-chain/core"
)

type PeerClient struct {
	Peers []string
}

func NewPeerClient(peers []string) *PeerClient {
	return &PeerClient{Peers: peers}
}

func (pc *PeerClient) BroadcastBlock(block *core.Block) {
	for _, peer := range pc.Peers {
		url := fmt.Sprintf("%s/receive-block", peer)
		go func(peerURL string) {
			body, _ := json.Marshal(block)
			resp, err := http.Post(peerURL, "application/json", bytes.NewReader(body))
			if err != nil {
				fmt.Printf("Failed to send block to %s: %v\n", peerURL, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Block sent to %s, response: %s\n", peerURL, resp.Status)
		}(url)
	}
}
