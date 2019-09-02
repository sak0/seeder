package repoer

import (
	"time"
	"github.com/sak0/go-harbor"
	"github.com/golang/glog"
	"encoding/json"
	"github.com/sak0/seeder/pkg/utils"
	)

const (
	defaultWatchInterval 	= 10 * time.Second

	StatusActive 			= "active"
)

type NodeInfo struct {
	NodeName 		string		`json:"node_name"`
	NodeRole 		string		`json:"node_role"`
	AdvertiseAddr 	string		`json:"advertise_addr"`
	BindAddr 		string		`json:"bind_addr"`
	RepoAddr 		string 		`json:"repo_addr"`
	ImageCount 		int			`json:"image_count"`
	ChartCount 		int			`json:"chart_count"`
	PullCount 		int			`json:"pull_count"`
	Status 			string		`json:"status"`
}

type ReporterInfo struct {
	NodeName 	string 							`json:"node_name"`
	NodeInfo    *NodeInfo 						`json:"node_info"`
	Repos		[]harbor.RepoRecord				`json:"repos"`
	Tags 		[]harbor.TagResp				`json:"tags"`
	Charts 		[]harbor.ChartRepoRecord		`json:"charts"`
	Versions 	[]harbor.ChartVersionRecord		`json:"versions"`
}

type RepoWatcher struct {
	stop 			chan interface{}
	watchInterval 	time.Duration
	client 			*harbor.Client
	info 			*ReporterInfo
}

func (w *RepoWatcher) Report() []byte {
	bytes, err := json.Marshal(w.info)
	if err != nil {
		glog.V(2).Infof("marshal report info failed: %v")
		return nil
	}
	return bytes
}

func (w *RepoWatcher) Run() {
	tick := time.NewTicker(w.watchInterval)

	for {
		select {
		case <-w.stop:
			return
		case <-tick.C:
			w.doLoop()
		}
	}
}

func (w *RepoWatcher) doLoop() {
	proOpts := harbor.ListProjectsOptions{
		Name : utils.DefaultProjectName,
	}
	projects, _, errs := w.client.Projects.ListProject(&proOpts)
	if len(errs) > 0 {
		glog.V(2).Infof("list repository failed: %v", errs)
		return
	}

	repoOpts := harbor.ListRepositoriesOption{
		ProjectId: projects[0].ProjectID,
	}
	repos, _, errs := w.client.Repositories.ListRepository(&repoOpts)
	if len(errs) > 0 {
		glog.V(2).Infof("list repository failed: %v", errs)
		return
	}
	w.info.Repos = repos

	var totalTags []harbor.TagResp
	for _, repo := range repos {
		tags, _, errs := w.client.Repositories.ListRepositoryTags(repo.Name)
		if len(errs) > 0 {
			glog.V(2).Infof("list repository tags failed: %v", errs)
			continue
		}
		totalTags = append(totalTags, tags...)
		for _, tag := range tags {
			glog.V(5).Infof("[scan-image] %s:%s", repo.Name, tag.Name)
		}
	}
	w.info.Tags = totalTags

	charts, _, errs := w.client.ChartRepos.ListChartRepositories(utils.DefaultProjectName)
	if len(errs) > 0 {
		glog.V(2).Infof("list chartRepos failed: %v", errs)
		return
	}
	w.info.Charts = charts

	var totalVersions []harbor.ChartVersionRecord
	for _, chart := range charts {
		versions, _, errs := w.client.ChartRepos.ListChartVersions(utils.DefaultProjectName, chart.Name)
		if len(errs) > 0 {
			glog.V(2).Infof("list chartRepos %s version failed: %v", chart.Name, errs)
			continue
		}
		totalVersions = append(totalVersions, versions...)
		for _, version := range versions {
			glog.V(5).Infof("[scan-chart] %s:%s", chart.Name, version.Version)
		}
	}
	w.info.Versions = totalVersions
}

func NewRepoWatcher(nodeName, nodeRole, repoAddr, advAddr, bindAddr string,
			stopCh chan interface{}) (*RepoWatcher, error){
	harborClient := harbor.NewClient(nil, repoAddr,utils.HarborUser,utils.HarborPass)
	opt := harbor.ListProjectsOptions{Name: utils.DefaultProjectName}
	projects, _, errs := harborClient.Projects.ListProject(&opt)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	glog.V(2).Infof("%+v", projects)

	return &RepoWatcher{
		client: harborClient,
		stop:stopCh,
		watchInterval:defaultWatchInterval,
		info:&ReporterInfo{
			NodeName: nodeName,
			NodeInfo: &NodeInfo{
				NodeName:nodeName,
				NodeRole:nodeRole,
				RepoAddr:repoAddr,
				AdvertiseAddr:advAddr,
				BindAddr:bindAddr,
				Status:StatusActive,
			},
		},
	}, nil
}