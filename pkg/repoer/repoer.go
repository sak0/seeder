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
	glog.V(2).Infof("%+v", repos)
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