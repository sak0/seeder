package keeper

import (
	"testing"
	"github.com/sak0/seeder/pkg/repoer"
	"github.com/sak0/seeder/pkg/utils"
)

var KeepInfo repoer.ReporterInfo

func TestFormatNode(t *testing.T) {
	nodeInfo := &repoer.NodeInfo{
		NodeName:"test-node",
		NodeRole:"master",
		AdvertiseAddr:utils.MustGetMyIpAddr(),
	}

	KeepInfo = repoer.ReporterInfo{
		NodeName:"test-node",
		NodeInfo:nodeInfo,
	}

	_, err := formatNode(KeepInfo)
	if err != nil {
		t.Fatalf("formatNode failed: %v", err)
	}
}