package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/sak0/memberlist"
	"github.com/sak0/seeder/pkg/utils"
	"github.com/sak0/seeder/pkg/repoer"
)

const (
	defaultLoopInterval = 10 * time.Second
)

type MyDelegate struct {
	syncer 	ClusterSyncer
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
	return []byte("node meta")
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
	glog.V(5).Infof("NotifyMsg: %v", string(msg))
	var repoInfo repoer.ReporterInfo
	if err := json.Unmarshal(msg, &repoInfo); err != nil {
		glog.V(2).Infof("receive invalid msg: %v", err)
		return
	}
	d.unpdateInfo(d.syncer, repoInfo)
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	//return [][]byte{[]byte("get broadcast")}
	return nil
}
func (d *MyDelegate) LocalState(join bool) []byte {
	return []byte("local state")
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
	glog.V(5).Infof("MergeRemoteState %s", buf)
}
func (d *MyDelegate) unpdateInfo(syncer ClusterSyncer, info repoer.ReporterInfo) {
	syncer.UpdateInfo(info)
}

type ClusterSyncer interface {
	Run()
	RegisterReporter(watcher *repoer.RepoWatcher)
	UpdateInfo(repoer.ReporterInfo)
}

type SeederNode struct {
	Name 			string
	Addr 			string
	Role 			string
	Master 			string
	stop 			chan interface{}
	loopInterval 	time.Duration
	mList  			*memberlist.Memberlist
	watcher 		*repoer.RepoWatcher
	infoMap			map[string]repoer.ReporterInfo
}

func (n *SeederNode) UpdateInfo(info repoer.ReporterInfo) {
	n.infoMap[info.NodeRole] = info
}

func (n *SeederNode) RegisterReporter(watcher *repoer.RepoWatcher) {
	n.watcher = watcher
}

func (n *SeederNode) Run() {
	myDlg := &MyDelegate{
		syncer:n,
	}

	lanConfig := memberlist.DefaultLANConfig()
	lanConfig.Name = n.Name
	lanConfig.BindAddr = n.Addr
	lanConfig.Delegate = myDlg
	lanConfig.LogOutput = ioutil.Discard

	member, err := memberlist.Create(lanConfig)
	if err != nil {
		panic(err)
	}
	_, err = member.Join([]string{n.Master})
	if err != nil {
		panic(err)
	}
	n.mList = member
	glog.V(2).Infof("node %s join master %s succeed.", n.Name, n.Master)

	n.runSeederNode()
}

func (n *SeederNode) broadcastRepoInfo(info []byte) {
	glog.V(5).Infof("master %s broadcast repo info to all nodes", n.Name)
	for _, node := range n.mList.Members() {
		if strings.HasPrefix(node.Name, "master") {
			continue
		}
		n.mList.SendToTCP(node, info)
	}
}

func (n *SeederNode) doLoop() {
	for _, node := range n.mList.Members() {
		if strings.HasPrefix(node.Name, "master") {
			master := node
			if master.Name != n.Name {
				n.mList.SendToTCP(master, []byte(fmt.Sprintf("hello I'm %s", n.Name)))
			}
		}
	}

	glog.V(2).Infof("memberList: %v", n.mList.Members())
	received := n.watcher.Report()
	if received != nil {
		var reportInfo repoer.ReporterInfo
		err := json.Unmarshal(received, &reportInfo)
		if err != nil {
			glog.V(2).Infof("unmarshal reportInfo failed: %v", err)
			return
		}
		glog.V(5).Infof("receive local report: %v", reportInfo)
		if n.Role == "master" {
			n.broadcastRepoInfo(received)
		}
	}
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
		masterAddr = utils.MustGetMyIpAddr()
	}

	return &SeederNode{
		Name:nodeName,
		Addr: utils.MustGetMyIpAddr(),
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