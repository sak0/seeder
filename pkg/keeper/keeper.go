package keeper

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/sak0/seeder/pkg/cluster"
	"github.com/sak0/seeder/pkg/repoer"
	"github.com/sak0/seeder/models"
	"github.com/sak0/seeder/pkg/utils"
)

const (
	defaultKeepInterval	= 30 * time.Second

	verifyStatusTrue 	= "verified"
	verifyStatusFalse 	= "unverified"
	verifyStatusUnknown	= "unknown"
)

func judgementVerifyStatus() string {
	if utils.MyRole == "master" {
		return verifyStatusTrue
	}
	return verifyStatusFalse
}

func formatNode(keepInfo repoer.ReporterInfo) (*models.SeederNode, error) {
	if keepInfo.NodeInfo == nil {
		return nil, fmt.Errorf("invalid report info: %s", keepInfo.NodeName)
	}

	return &models.SeederNode{
		ClusterName:keepInfo.NodeInfo.NodeName,
		AdvertiseAddr:keepInfo.NodeInfo.AdvertiseAddr,
		BindAddr:keepInfo.NodeInfo.BindAddr,
		RepoAddr:keepInfo.NodeInfo.RepoAddr,
		Role:keepInfo.NodeInfo.NodeRole,
		ImageCount:keepInfo.NodeInfo.ImageCount,
		ChartCount:keepInfo.NodeInfo.ChartCount,
		PullCount:keepInfo.NodeInfo.PullCount,
		Status:keepInfo.NodeInfo.Status,
	}, nil
}

func formatVersions(keepInfo repoer.ReporterInfo, cached bool, verifyStatus string) []*models.ChartVersion {
	var versions []*models.ChartVersion
	for _, chartVersion := range keepInfo.Versions {
		version := &models.ChartVersion {
			Name:chartVersion.Name,
			Version:chartVersion.Version,
			Description:chartVersion.Description,
			AppVersion:chartVersion.AppVersion,
			Url:chartVersion.Urls[0],
			Digest:chartVersion.Digest,
			VerifyStatus:verifyStatus,
			Cached:cached,
			CreationTime:chartVersion.CreationTime,
			UpdateTime:chartVersion.UpdateTime,
		}
		versions = append(versions, version)
	}

	return versions
}

func formatCharts(keepInfo repoer.ReporterInfo, cached bool, verifyStatus string) []*models.ChartRepo {
	var charts []*models.ChartRepo
	for _, chartRepo := range keepInfo.Charts {
		chart := &models.ChartRepo{
			OwnerNode:"master",
			Name:chartRepo.Name,
			VersionCount:chartRepo.TotalVersions,
			LatestVersion:chartRepo.LatestVersion,
			Icon:chartRepo.Icon,
			Home:chartRepo.Home,
			VerifyStatus:verifyStatus,
			Cached:cached,
			CreationTime:chartRepo.CreationTime,
			UpdateTime:chartRepo.UpdateTime,
		}
		charts = append(charts, chart)
	}

	return charts
}

func formatTags(keepInfo repoer.ReporterInfo) []*models.RepositoryTag {
	var tags []*models.RepositoryTag
	for _, repoTag := range keepInfo.Tags {
		tag := &models.RepositoryTag{
			Digest:repoTag.Digest,
			TagName:repoTag.Name,
			Size:repoTag.Size,
			Architecture:repoTag.Architecture,
			OS:repoTag.OS,
			DockerVersion:repoTag.DockerVersion,
			Author:repoTag.Author,
			VerifyStatus:verifyStatusTrue,
			Cached:true,
		}
		tags = append(tags, tag)
	}

	return tags
}

func formatRepos(keepInfo repoer.ReporterInfo) []*models.Repository {
	var repos []*models.Repository
	for _, infoRepo := range keepInfo.Repos {
		repo := &models.Repository{
			OwnerNode: keepInfo.NodeInfo.NodeName,
			Name: infoRepo.Name,
			Description: infoRepo.Description,
			PullCount:infoRepo.PullCount,
			StarCount:infoRepo.StarCount,
			VerifyStatus:verifyStatusTrue,
			Cached:true,
		}
		repos = append(repos, repo)
	}
	return repos
}

func diffNodes(remote *models.SeederNode, local []*models.SeederNode) ([]*models.SeederNode, []*models.SeederNode, []*models.SeederNode) {
	var addNodes []*models.SeederNode

	found := false
	for _, localNode := range local {
		if remote.ClusterName == localNode.ClusterName {
			found = true
			continue
		}
	}
	if !found {
		addNodes = append(addNodes, remote)
	}

	return addNodes, nil, nil
}

func diffTags(remote, local []*models.RepositoryTag) ([]*models.RepositoryTag, []*models.RepositoryTag, []*models.RepositoryTag) {
	var addTags []*models.RepositoryTag

Loop:
	for _, remoteTag := range remote {
		found := false
		for _, localTag := range local {
			if remoteTag.TagName == localTag.TagName {
				found = true
				continue Loop
			}
		}
		if !found {
			addTags = append(addTags, remoteTag)
		}
	}

	return addTags, nil, nil
}

func diffVersions(remote, local []*models.ChartVersion) ([]*models.ChartVersion, []*models.ChartVersion, []*models.ChartVersion) {
	var addVersions []*models.ChartVersion
Loop:
	for _, remoteVersion := range remote {
		found := false
		for _, localVersion := range local {
			if remoteVersion.Name == localVersion.Name {
				found = true
				continue Loop
			}
		}
		if !found {
			addVersions = append(addVersions, remoteVersion)
		}
	}

	return addVersions, nil, nil
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

func (k *LocalKeeper) getClusterInfo() map[string]repoer.ReporterInfo {
	return k.reporter.GetInfoMap()
}

func (k *LocalKeeper) getKeepInfo() (repoer.ReporterInfo, error){
	clusterInfo := k.getClusterInfo()

	keepInfo, ok := clusterInfo[k.name]
	if !ok {
		return repoer.ReporterInfo{}, fmt.Errorf(fmt.Sprintf("miss %s info from reporter", k.name))
	}
	return keepInfo, nil
}

func (k *LocalKeeper) getMasterInfo() (repoer.ReporterInfo, error) {
	clusterInfo := k.getClusterInfo()

	masters, err := models.GetNodesByRole("master")
	if err != nil {
		return repoer.ReporterInfo{}, err
	}

	masterInfo, ok := clusterInfo[masters[0].ClusterName]
	if !ok {
		return repoer.ReporterInfo{}, fmt.Errorf(fmt.Sprintf("miss %s info from master", masters[0].ClusterName))
	}
	return masterInfo, nil
}

func (k *LocalKeeper) getLocalNodes() ([]*models.SeederNode, error) {
	nodes, _, err := models.GetSeederNodes(0, 0)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (k *LocalKeeper) getLocalUnCachedVersions() ([]*models.ChartVersion, error) {
	versions, _, err := models.GetUnCachedVersions(0, 0, false)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (k *LocalKeeper) getLocalVersions() ([]*models.ChartVersion, error) {
	versions, _, err := models.GetAllVersions(0, 0)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (k *LocalKeeper) getLocalCharts() ([]*models.ChartRepo, error) {
	charts, _, err := models.GetAllCharts(0, 0, "")
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

func (k *LocalKeeper) addNode(node *models.SeederNode) {
	glog.V(2).Infof("ADD NODE: %v", node)
	if err := models.CreateNode(node); err != nil {
		glog.V(2).Infof("add node failed: %v", err)
	}
}

func (k *LocalKeeper) addTag(tag *models.RepositoryTag) {
	glog.V(2).Infof("ADD TAG: %v", tag)
	if err := models.CreateTag(tag); err != nil {
		glog.V(2).Infof("add tag failed: %v", err)
	}
}

func (k *LocalKeeper) addVersion(version *models.ChartVersion) {
	glog.V(2).Infof("ADD VERSION: %v", version)
	if err := models.CreateVersion(version); err != nil {
		glog.V(2).Infof("add version failed: %v", err)
	}
}

func (k *LocalKeeper) markVersionCached(chartName, version string) {
	glog.V(2).Infof("UPDATE VERSION CACHED: %s/%s", chartName, version)
	if err := models.UpdateVersionCached(chartName, version); err != nil {
		glog.V(2).Infof("update version failed: %v", err)
	}
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

func (k *LocalKeeper) syncNode(keepInfo repoer.ReporterInfo) {
	remoteNode, err := formatNode(keepInfo)
	if err != nil {
		glog.V(2).Infof("get remote nodes failed: %v", err)
		return
	}
	localNodes, err := k.getLocalNodes()
	if err != nil {
		glog.V(2).Infof("get local nodes failed: %v", err)
		return
	}
	glog.V(2).Infof("[remoteNode] %v", remoteNode)
	glog.V(2).Infof("[localNodes] %v", localNodes)
	nodeAdds, _, _ := diffNodes(remoteNode, localNodes)
	for _, nodeAdd := range nodeAdds {
		k.addNode(nodeAdd)
	}
}

func (k *LocalKeeper) syncRepos(keepInfo repoer.ReporterInfo) {
	remoteRepos := formatRepos(keepInfo)
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

func (k *LocalKeeper) syncTags(keepInfo repoer.ReporterInfo) {
	remoteTags := formatTags(keepInfo)
	localTags, err := k.getLocalTags()
	if err != nil {
		glog.V(2).Infof("get local tags failed: %v", err)
		return
	}
	glog.V(5).Infof("[remoteTags] %v", remoteTags)
	glog.V(5).Infof("[localTags] %v", localTags)
	tagsAdd, _, _ := diffTags(remoteTags, localTags)
	for _, tagAdd := range tagsAdd {
		k.addTag(tagAdd)
	}
}

func (k *LocalKeeper) syncCharts(keepInfo, masterInfo repoer.ReporterInfo) {
	k.syncLocalCharts(keepInfo)
	k.syncMasterCharts(masterInfo)
}

func (k *LocalKeeper) syncLocalCharts(keepInfo repoer.ReporterInfo) {
	remoteCharts := formatCharts(keepInfo, true, judgementVerifyStatus())
	localCharts, err := k.getLocalCharts()
	if err != nil {
		glog.V(2).Infof("get local charts failed: %v", err)
		return
	}
	glog.V(5).Infof("[remoteCharts] %v", remoteCharts)
	glog.V(5).Infof("[localCharts] %v", localCharts)
	chartsAdd, _, _ := diffCharts(remoteCharts, localCharts)
	for _, chartAdd := range chartsAdd {
		k.addChart(chartAdd)
	}
}

func (k *LocalKeeper) syncMasterCharts(masterInfo repoer.ReporterInfo) {

}

func (k *LocalKeeper) syncVersions(keepInfo, masterInfo repoer.ReporterInfo) {
	k.syncLocalVersions(keepInfo)
	k.syncMasterVersions(masterInfo)
}

func (k *LocalKeeper) syncMasterVersions(masterInfo repoer.ReporterInfo) {
	masterVersions := formatVersions(masterInfo, false, verifyStatusTrue)
	localVersions, err := k.getLocalVersions()
	if err != nil {
		glog.V(2).Infof("get local charts failed: %v", err)
		return
	}

	glog.V(5).Infof("[masterVersions] %v", masterVersions)
	glog.V(5).Infof("[localVersions] %v", localVersions)
	versionsAdd, _, _ := diffVersions(masterVersions, localVersions)
	for _, versionAdd := range versionsAdd {
		k.addVersion(versionAdd)
	}
}

func (k *LocalKeeper) syncLocalVersions(keepInfo repoer.ReporterInfo) {
	remoteVersions := formatVersions(keepInfo, true, judgementVerifyStatus())
	localVersions, err := k.getLocalVersions()
	if err != nil {
		glog.V(2).Infof("get local charts failed: %v", err)
		return
	}
	glog.V(5).Infof("[remoteVersions] %v", remoteVersions)
	glog.V(5).Infof("[localVersions] %v", localVersions)
	versionsAdd, _, _ := diffVersions(remoteVersions, localVersions)
	for _, versionAdd := range versionsAdd {
		k.addVersion(versionAdd)
	}

	if utils.MyRole == "follower" {
		localUnCachedVersions, err := k.getLocalUnCachedVersions()
		if err != nil {
			glog.V(2).Infof("get local unCached versions failed: %v", err)
			return
		}
		glog.V(5).Infof("[remoteVersions] %v", remoteVersions)
		glog.V(5).Infof("[localUnCachedVersions] %v", localUnCachedVersions)
		for _, localUnCachedVersion := range localUnCachedVersions {
			for _, remoteVersion := range remoteVersions {
				if localUnCachedVersion.Version == remoteVersion.Version &&
					localUnCachedVersion.Name == remoteVersion.Name {
					k.markVersionCached(localUnCachedVersion.Name, localUnCachedVersion.Version)
				}
			}
		}
	}
}

func (k *LocalKeeper) syncUnCachedCharts(keepInfo, masterInfo repoer.ReporterInfo) {

}

func (k *LocalKeeper) syncUnCachedVersions(keepInfo, masterInfo repoer.ReporterInfo) {

}

func (k *LocalKeeper) doSync() {
	clusterInfo := k.getClusterInfo()
	for _, info := range clusterInfo {
		k.syncNode(info)
	}

	keepInfo, err := k.getKeepInfo()
	if err != nil {
		glog.V(2).Infof("get keep info failed: %v", err)
		return
	}
	glog.V(5).Infof("keeper sync: %v", keepInfo)

	masterInfo, err := k.getMasterInfo()
	if err != nil {
		glog.V(2).Infof("get master info failed: %v", err)
		return
	}

	{
		k.syncRepos(keepInfo)
		k.syncTags(keepInfo)
		k.syncCharts(keepInfo, masterInfo)
		k.syncVersions(keepInfo, masterInfo)
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