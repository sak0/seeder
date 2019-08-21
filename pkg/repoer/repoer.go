package repoer

import (
	"time"
	"github.com/sak0/go-harbor"
	"github.com/golang/glog"
)

const (
	defaultWatchInterval 	= 10 * time.Second
	defaultProjectName 		= "edge-cloud"
)

type RepoWatcher struct {
	stop 			chan interface{}
	watchInterval 	time.Duration
	client 			*harbor.Client
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
	repoOpts := harbor.ListRepositoriesOption{
		ProjectId: 3,
	}
	repos, _, errs := w.client.Repositories.ListRepository(&repoOpts)
	if len(errs) > 0 {
		glog.V(2).Infof("list repository failed: %v", errs)
		return
	}
	for _, repo := range repos {
		tags, _, errs := w.client.Repositories.ListRepositoryTags(repo.Name)
		if len(errs) > 0 {
			glog.V(2).Infof("list repository tags failed: %v", errs)
			continue
		}
		for _, tag := range tags {
			glog.V(2).Infof("[IMAGE] %s:%s", repo.Name, tag.Name)
		}
	}

	charts, _, errs := w.client.ChartRepos.ListChartRepositories(defaultProjectName)
	if len(errs) > 0 {
		glog.V(2).Infof("list chartRepos failed: %v", errs)
		return
	}
	for _, chart := range charts {
		versions, _, errs := w.client.ChartRepos.ListChartVersions(defaultProjectName, chart.Name)
		if len(errs) > 0 {
			glog.V(2).Infof("list chartRepos %s version failed: %v", chart.Name, errs)
			continue
		}
		for _, version := range versions {
			glog.V(2).Infof("[CHART] %s:%s", chart.Name, version.Version)
		}
	}
}

func NewRepoWatcher(repoAddr string, stopCh chan interface{}) (*RepoWatcher, error){
	harborClient := harbor.NewClient(nil, repoAddr,"admin","Harbor12345")
	opt := harbor.ListProjectsOptions{Name: defaultProjectName}
	projects, _, errs := harborClient.Projects.ListProject(&opt)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	glog.V(2).Infof("%+v", projects)

	return &RepoWatcher{
		client: harborClient,
		stop:stopCh,
		watchInterval:defaultWatchInterval,
	}, nil
}