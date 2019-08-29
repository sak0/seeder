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

type NodeMeta struct {
	Name 	string
	Role 	string
	Addr 	string
}

type MyDelegate struct {
	syncer 	ClusterSyncer
	meta 	NodeMeta
}

func (d *MyDelegate) NodeMeta(limit int) []byte {
	bytes, err := json.Marshal(&d.meta)
	if err != nil {
		glog.V(2).Infof("get node meta failed: %v", err)
		return nil
	}
	return bytes
}

func (d *MyDelegate) NotifyMsg(msg []byte) {
	glog.V(5).Infof("NotifyMsg: %v", string(msg))
	var repoInfo repoer.ReporterInfo
	if err := json.Unmarshal(msg, &repoInfo); err != nil {
		glog.V(2).Infof("receive invalid msg: %v", err)
		return
	}
	d.updateInfo(d.syncer, repoInfo)
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

func (d *MyDelegate) updateInfo(syncer ClusterSyncer, info repoer.ReporterInfo) {
	syncer.UpdateInfo(info)
}

type ClusterSyncer interface {
	Run()
	RegisterReporter(watcher *repoer.RepoWatcher)
	UpdateInfo(repoer.ReporterInfo)
	GetInfoMap()map[string]repoer.ReporterInfo
	GetNodes()map[string]string
}

type SeederNode struct {
	Name 			string
	Addr 			string
	Role 			string
	Master 			string
	RepoAddr 		string
	stop 			chan interface{}
	loopInterval 	time.Duration
	mList  			*memberlist.Memberlist
	watcher 		*repoer.RepoWatcher
	nodes           map[string]string
	infoMap			map[string]repoer.ReporterInfo
}

func (n *SeederNode) GetInfoMap() map[string]repoer.ReporterInfo {
	return n.infoMap
}

func (n *SeederNode) GetNodes() map[string]string {
	return n.nodes
}

func (n *SeederNode) UpdateInfo(info repoer.ReporterInfo) {
	n.infoMap[info.NodeName] = info
}

func (n *SeederNode) RegisterReporter(watcher *repoer.RepoWatcher) {
	n.watcher = watcher
}

func (n *SeederNode) Run() {
	myDlg := &MyDelegate{
		syncer:n,
		meta:NodeMeta{
			Name: n.Name,
			Role: n.Role,
			Addr: n.Addr,
		},
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
	nodes := make(map[string]string)
	for _, node := range n.mList.Members() {
		if strings.HasPrefix(node.Name, "master") {
			nodes[node.Name] = "master"
			master := node
			if master.Name != n.Name {
				n.mList.SendToTCP(master, []byte(fmt.Sprintf("hello I'm %s", n.Name)))
			}
		} else {
			nodes[node.Name] = "follower"
		}
	}
	n.nodes = nodes

	members := n.mList.Members()
	for _, member := range members {
		var meta NodeMeta
		err := json.Unmarshal(member.Meta, &meta)
		if err != nil {
			glog.V(2).Infof("can not marshal node %s meta", member.Name)
			continue
		}
		glog.V(2).Infof("[%s-%s] %v", member.Name, member.Addr, meta)
	}

	received := n.watcher.Report()
	if received != nil {
		var reportInfo repoer.ReporterInfo
		err := json.Unmarshal(received, &reportInfo)
		if err != nil {
			glog.V(2).Infof("unmarshal reportInfo failed: %v", err)
			return
		}
		glog.V(5).Infof("receive local report: %v", reportInfo)
		//if n.Role == "master" {
		//	n.broadcastRepoInfo(received)
		//}
		n.broadcastRepoInfo(received)
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

func newSeederNode(role, masterAddr, nodeName, repoAddr string, stopCh chan interface{}) *SeederNode {
	if masterAddr == "" {
		masterAddr = utils.MustGetMyIpAddr()
	}

	return &SeederNode{
		Name:nodeName,
		Addr: utils.MustGetMyIpAddr(),
		Role:role,
		RepoAddr:repoAddr,
		Master:masterAddr,
		stop: stopCh,
		loopInterval: defaultLoopInterval,
		infoMap:make(map[string]repoer.ReporterInfo),
		nodes:make(map[string]string),
	}
}

func NewClusterSyncer(role, masterAddr, nodeName, repoAddr, syncMode string, stopCh chan interface{}) ClusterSyncer {
	switch syncMode {
	case "gossip":
		syncer := newSeederNode(role, masterAddr, nodeName, repoAddr, stopCh)
		return syncer
	default:
		panic(fmt.Sprintf("unsupport syncer mode: %s", syncMode))
	}

	return nil
}