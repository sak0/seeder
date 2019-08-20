package cluster

import (
	"fmt"
	"git.dev.yonghui.cn/ieyes/pajero/utils"
	"github.com/golang/glog"
	"github.com/sak0/memberlist"
	"time"
)

const (
	defaultLoopInterval = 10 * time.Second
)

type MyDelegate struct {
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
	return []byte("node meta")
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
	glog.V(5).Infof("NotifyMsg: %v", string(msg))
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{[]byte("brdcast")}
}
func (d *MyDelegate) LocalState(join bool) []byte {
	return []byte("local state")
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
	glog.V(5).Infof("MergeRemoteState %v", buf)
}

type ClusterSyncer interface {
	Run()
}

type SeederNode struct {
	Name 			string
	Addr 			string
	Role 			string
	Master 			string
	stop 			chan interface{}
	loopInterval 	time.Duration
	mList  			*memberlist.Memberlist
}
func (n *SeederNode) Run() {
	myDlg := &MyDelegate{}

	lanConfig := memberlist.DefaultLANConfig()
	lanConfig.Name = n.Name
	lanConfig.BindAddr = n.Addr
	lanConfig.Delegate = myDlg

	member, err := memberlist.Create(lanConfig)
	if err != nil {
		panic(err)
	}
	_, err = member.Join([]string{n.Master})
	if err != nil {
		panic(err)
	}
	n.mList = member

	n.runSeederNode()
}

func (n *SeederNode) doLoop() {
	glog.V(2).Infof("%v", n.mList.Members())
}

func (n *SeederNode) runSeederNode() {
	tick := time.NewTicker(n.loopInterval)

	for {
		select {
		case <-n.stop:
			return
		case <-tick.C:
			n.doLoop()
		}
	}
}

func newSeederNode(role, masterAddr, nodeName string, stopCh chan interface{}) *SeederNode {
	if masterAddr == "" {
		masterAddr = utils.GetMyIPAddr()
	}

	return &SeederNode{
		Name:nodeName,
		Addr: utils.GetMyIPAddr(),
		Role:role,
		Master:masterAddr,
		stop: stopCh,
		loopInterval: defaultLoopInterval,
	}
}

func NewClusterSyncer(role, masterAddr, nodeName, syncMode string, stopCh chan interface{}) ClusterSyncer {
	switch syncMode {
	case "gossip":
		syncer := newSeederNode(role, masterAddr, nodeName, stopCh)
		return syncer
	default:
		panic(fmt.Sprintf("unsupport syncer mode: %s", syncMode))
	}

	return nil
}