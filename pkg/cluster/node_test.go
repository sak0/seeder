package cluster

import (
	"testing"
	"github.com/sak0/seeder/pkg/utils"
)

func TestNewClusterSyncer(t *testing.T) {
	stopCh := make(chan interface{})
	_ := NewClusterSyncer("follower", utils.MustGetMyIpAddr(), "node-test",
		"http://172.16.24.103", "gossip", stopCh)
}

func TestSeederNode_Run(t *testing.T) {
	stopCh := make(chan interface{})
	sync := NewClusterSyncer("follower", utils.MustGetMyIpAddr(), "node-test",
		"http://172.16.24.103", "gossip", stopCh)
	sync.Run()
}