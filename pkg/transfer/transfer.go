package transfer

import (
	"strings"
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/sak0/go-harbor"
	"github.com/sak0/seeder/pkg/repoer"
)

type Transfer struct {
	srcAddr 		string
	dstAddr 		string
	SrcRepo 		*harbor.Client
	DstRepo 		*harbor.Client
	client 			*http.Client
}

type label struct {
	Name string `json:"name"`
}

type chartVersion struct {
	Version string   `json:"version"`
	Labels  []*label `json:"labels"`
}

type chartVersionDetail struct {
	Metadata *chartVersionMetadata `json:"metadata"`
}

type chartVersionMetadata struct {
	URLs []string `json:"urls"`
}

func parseChartName(name string) (string, string, error) {
	strs := strings.Split(name, "/")
	if len(strs) == 2 && len(strs[0]) > 0 && len(strs[1]) > 0 {
		return strs[0], strs[1], nil
	}
	return "", "", fmt.Errorf("invalid chart name format: %s", name)
}

func (t *Transfer) getChartInfo(name, version string) (*chartVersionDetail, error) {
	project, name, err := parseChartName(name)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/api/chartrepo/%s/charts/%s/%s", t.srcAddr, project, name, version)
	info := &chartVersionDetail{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (t *Transfer) DownloadChart(chartName string, chartVersion string) (io.ReadCloser, error){
	info, err := t.getChartInfo(chartName, chartVersion)
	if err != nil {
		return nil, err
	}
	if info.Metadata == nil || len(info.Metadata.URLs) == 0 || len(info.Metadata.URLs[0]) == 0 {
		return nil, fmt.Errorf("cannot got the download url for chart %s:%s", chartName, chartVersion)
	}
	url := strings.ToLower(info.Metadata.URLs[0])
	// relative URL
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		project, _, err := parseChartName(chartName)
		if err != nil {
			return nil, err
		}
		url = fmt.Sprintf("%s/chartrepo/%s/%s", t.srcAddr, project, url)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to download the chart %s: %d %s", req.URL.String(), resp.StatusCode, string(body))
	}
	return resp.Body, nil
}


func mustHarborClient(repoAddr string)(*harbor.Client, error) {
	client := harbor.NewClient(nil, repoAddr, "admin", "Harbor12345")
	opt := harbor.ListProjectsOptions{Name: repoer.DefaultProjectName}
	_, _, errs := client.Projects.ListProject(&opt)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return client, nil
}

func NewTransfer(srcAddr, dstAddr string) (*Transfer, error) {
	sc, err := mustHarborClient(srcAddr)
	if err != nil {
		return nil, err
	}
	dc, err := mustHarborClient(dstAddr)
	if err != nil {
		return nil, err
	}

	return &Transfer{
		srcAddr:srcAddr,
		dstAddr:dstAddr,
		SrcRepo:sc,
		DstRepo:dc,
	}, nil
}