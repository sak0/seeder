package utils

import (
	"testing"
	"github.com/golang/glog"
)

func TestMustGetMyIpAddr(t *testing.T) {
	addr, err := GetMyIpAddr()
	if err != nil {
		t.Fatalf("test get ip addr failed: %v", err)
	}
	glog.V(2).Infof("test get ip addr: %s", addr)
}