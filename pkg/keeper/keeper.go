package keeper

import (
	"github.com/sak0/seeder/pkg/cluster"
	"time"
	"github.com/sak0/seeder/pkg/repoer"
	"fmt"
	"github.com/golang/glog"
	"github.com/sak0/seeder/models"
)

const (
	defaultKeepInterval	= 30 * time.Second
)

func formatRepos(masterInfo repoer.ReporterInfo) []*models.Repository {
	var repos []*models.Repository
	for _, infoRepo := range masterInfo.Repos {
		repo := &models.Repository{
			OwnerNode: masterInfo.NodeName,
			Name: infoRepo.Name,
			Description: infoRepo.Description,
			PullCount:infoRepo.PullCount,
			StarCount:infoRepo.StarCount,
			VerifyStatus:"verified",
			Cached:true,
		}
		repos = append(repos, repo)
	}
	return repos
}

func diffRepos(remote, local []*models.Repository) ([]*models.Repository, []*models.Repository, []*models.Repository) {
	var addRepos []*models.Repository

Loop:
	for _, remoteRepo := range remote {
		found := false
		for _, localRepo := range local {
			if remoteRepo.Name == localRepo.Name {
				found = true
				continue Loop
			}
		}
		if !found {
			addRepos = append(addRepos, remoteRepo)
		}
	}

	return addRepos, nil, nil
}

type LocalKeeper struct {
	name 		string
	role 		string
	master 		string
	stop		chan interface{}
	reporter 	cluster.ClusterSyncer
	interval    time.Duration
}

func (k *LocalKeeper) getMasterInfo() (repoer.ReporterInfo, error){
	var masterName string
	nodes := k.reporter.GetNodes()
	for name, role := range nodes {
		if role == "master" {
			masterName = name
		}
	}

	clusterInfo := k.reporter.GetInfoMap()
	masterInfo, ok := clusterInfo[masterName]
	if !ok {
		return repoer.ReporterInfo{}, fmt.Errorf(fmt.Sprintf("miss master info from reporter"))
	}
	return masterInfo, nil
}

func (k *LocalKeeper) getLocalRepos() ([]*models.Repository, error) {
	repos, _, err := models.GetAllRepos(0, 0)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (k *LocalKeeper) addRepo(repository *models.Repository) {
	glog.V(2).Infof("ADD REPO: %v", repository)
	if err := models.CreateRepo(repository); err != nil {
		glog.V(2).Infof("add repo failed: %v", err)
	}
}

func (k *LocalKeeper) syncRepos() {
	masterInfo, err := k.getMasterInfo()
	if err != nil {
		glog.V(2).Infof("%v", err)
		return
	}
	remoteRepos := formatRepos(masterInfo)
	localRepos, err := k.getLocalRepos()
	if err != nil {
		glog.V(2).Infof("get local repo failed: %v", err)
		return
	}
	glog.V(2).Infof("[remoteRepos] %v", remoteRepos)
	glog.V(2).Infof("[localRepos] %v", localRepos)
	reposAdd, _, _ := diffRepos(remoteRepos, localRepos)
	for _, repoAdd := range reposAdd {
		k.addRepo(repoAdd)
	}
}

func (k *LocalKeeper) doSync() {
	k.syncRepos()
}

func (k *LocalKeeper) Run() {
	tick := time.NewTicker(k.interval)
	defer tick.Stop()

	for {
		select {
		case <-k.stop:
			return
		case <-tick.C:
			k.doSync()
		}
	}
}

func (k *LocalKeeper) RegisterReporter(reporter cluster.ClusterSyncer) {
	k.reporter = reporter
}

func NewLocalKeeper(role, master, myName string, stop chan interface{}) *LocalKeeper {
	return &LocalKeeper{
		name:myName,
		role:role,
		master:master,
		stop:stop,
		interval:defaultKeepInterval,
	}
}