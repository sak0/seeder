package keeper

import (
	"fmt"
	"time"

	"github.com/sak0/seeder/pkg/cluster"
	"github.com/sak0/seeder/pkg/repoer"
	"github.com/golang/glog"
	"github.com/sak0/seeder/models"
)

const (
	defaultKeepInterval	= 30 * time.Second
)

func formatCharts(masterInfo repoer.ReporterInfo) []*models.ChartRepo {
	var charts []*models.ChartRepo
	for _, chartRepo := range masterInfo.Charts {
		chart := &models.ChartRepo{
			OwnerNode:"master",
			Name:chartRepo.Name,
			VersionCount:chartRepo.TotalVersions,
			LatestVersion:chartRepo.LatestVersion,
			Icon:chartRepo.Icon,
			Home:chartRepo.Home,
			VerifyStatus:"verified",
			Cached:true,
		}
		charts = append(charts, chart)
	}

	return charts
}

func formatTags(masterInfo repoer.ReporterInfo) []*models.RepositoryTag {
	var tags []*models.RepositoryTag
	for _, repoTag := range masterInfo.Tags {
		tag := &models.RepositoryTag{
			Digest:repoTag.Digest,
			TagName:repoTag.Name,
			Size:repoTag.Size,
			Architecture:repoTag.Architecture,
			OS:repoTag.OS,
			DockerVersion:repoTag.DockerVersion,
			Author:repoTag.Author,
			VerifyStatus:"verified",
			Cached:true,
		}
		tags = append(tags, tag)
	}

	return tags
}

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

func diffCharts(remote, local []*models.ChartRepo) ([]*models.ChartRepo, []*models.ChartRepo, []*models.ChartRepo) {
	var addCharts []*models.ChartRepo

Loop:
	for _, remoteChart := range remote {
		found := false
		for _, localChart := range local {
			if remoteChart.Name == localChart.Name {
				found = true
				continue Loop
			}
		}
		if !found {
			addCharts = append(addCharts, remoteChart)
		}
	}

	return addCharts, nil, nil
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

func (k *LocalKeeper) getLocalCharts() ([]*models.ChartRepo, error) {
	charts, _, err := models.GetAllCharts(0, 0)
	if err != nil {
		return nil, err
	}
	return charts, nil
}

func (k *LocalKeeper) getLocalTags() ([]*models.RepositoryTag, error) {
	tags, _, err := models.GetAllTags(0, 0)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (k *LocalKeeper) getLocalRepos() ([]*models.Repository, error) {
	repos, _, err := models.GetAllRepos(0, 0)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (k *LocalKeeper) addChart(chart *models.ChartRepo) {
	glog.V(2).Infof("ADD CHART: %v", chart)
	if err := models.CreateChart(chart); err != nil {
		glog.V(2).Infof("add chart failed: %v", err)
	}
}

func (k *LocalKeeper) addRepo(repository *models.Repository) {
	glog.V(2).Infof("ADD REPO: %v", repository)
	if err := models.CreateRepo(repository); err != nil {
		glog.V(2).Infof("add repo failed: %v", err)
	}
}

func (k *LocalKeeper) syncRepos(masterInfo repoer.ReporterInfo) {
	remoteRepos := formatRepos(masterInfo)
	localRepos, err := k.getLocalRepos()
	if err != nil {
		glog.V(2).Infof("get local repo failed: %v", err)
		return
	}
	glog.V(5).Infof("[remoteRepos] %v", remoteRepos)
	glog.V(5).Infof("[localRepos] %v", localRepos)
	reposAdd, _, _ := diffRepos(remoteRepos, localRepos)
	for _, repoAdd := range reposAdd {
		k.addRepo(repoAdd)
	}
}

func (k *LocalKeeper) syncTags(masterInfo repoer.ReporterInfo) {
	remoteTags := formatTags(masterInfo)
	localTags, err := k.getLocalTags()
	if err != nil {
		glog.V(2).Infof("get local tags failed: %v", err)
		return
	}
	glog.V(2).Infof("[remoteTags] %v", remoteTags)
	glog.V(2).Infof("[localTags] %v", localTags)
}

func (k *LocalKeeper) syncCharts(masterInfo repoer.ReporterInfo) {
	remoteCharts := formatCharts(masterInfo)
	localCharts, err := k.getLocalCharts()
	if err != nil {
		glog.V(2).Infof("get local charts failed: %v", err)
		return
	}
	glog.V(2).Infof("[remoteCharts] %v", remoteCharts)
	glog.V(2).Infof("[localCharts] %v", localCharts)
	chartsAdd, _, _ := diffCharts(remoteCharts, localCharts)
	for _, chartAdd := range chartsAdd {
		k.addChart(chartAdd)
	}
}

func (k *LocalKeeper) doSync() {
	masterInfo, err := k.getMasterInfo()
	if err != nil {
		glog.V(2).Infof("get master info failed: %v", err)
		return
	}

	{
		k.syncRepos(masterInfo)
		k.syncTags(masterInfo)
		k.syncCharts(masterInfo)
	}
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