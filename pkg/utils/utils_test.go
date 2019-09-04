package utils

import (
	"testing"
	"github.com/golang/glog"
	"os"
)

func TestMustGetMyIpAddr(t *testing.T) {
	addr, err := GetMyIpAddr()
	if err != nil {
		t.Fatalf("test get ip addr failed: %v", err)
	}
	glog.V(2).Infof("test get ip addr: %s", addr)
}

func TestHarborAuth(t *testing.T) {
	os.Setenv("HARBOR_USER", "admin")
	os.Setenv("HARBOR_PASSWD", "Harbor12345")

	err := HarborAuth()
	if err != nil {
		t.Fatalf("get harbor auth info failed: %v", err)
	}
}