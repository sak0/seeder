package transfer

import (
	"strings"
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"mime/multipart"

	"github.com/sak0/go-harbor"
		common_http "github.com/sak0/seeder/pkg/common/http"
	"github.com/golang/glog"
	"github.com/sak0/seeder/pkg/utils"
	"net/url"
)

type Transfer struct {
	srcAddr 		string
	dstAddr 		string
	SrcRepo 		*harbor.Client
	DstRepo 		*harbor.Client
	client   		*common_http.Client
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

	//req, err := http.NewRequest(http.MethodGet, url, nil)
	//if err != nil {
	//	return nil, err
	//}
	//resp, err := t.client.Do(req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//data, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal(data, &info)
	//if err != nil {
	//	return nil, err
	//}
	glog.V(2).Infof("getChart info from %s", url)
	err = t.client.Get(url, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (t *Transfer) downloadChart(chartName string, chartVersion string) (io.ReadCloser, error){
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

func (t *Transfer) uploadChart(name, version string, chart io.Reader) error {
	project, name, err := parseChartName(name)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("chart", name+".tgz")
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, chart); err != nil {
		return err
	}
	w.Close()

	url := fmt.Sprintf("%s/api/chartrepo/%s/charts", t.dstAddr, project)

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &common_http.Error{
			Code:    resp.StatusCode,
			Message: string(data),
		}
	}
	return nil
}

func (t *Transfer) Transfer(name, version string) error {
	fullName := utils.DefaultProjectName + "/" + name
	return t.copy(fullName, version, true)
}

func (t *Transfer) copy(name, version string, override bool) error {
	glog.V(2).Infof("copying %s:%s(source registry) to %s:%s(destination registry)...",
		name, version, name, version)

	// check the existence of the chart on the destination registry
	//exist, err := t.ChartExist(name, version)
	//if err != nil {
	//	glog.V(2).Infof("failed to check the existence of chart %s:%s on the destination registry: %v", name, version, err)
	//	return err
	//}
	//if exist {
	//	// the same name chart exists, but not allowed to override
	//	if !override {
	//		glog.V(2).Infof("the same name chart %s:%s exists on the destination registry, but the \"override\" is set to false, skip",
	//			name, version)
	//		return nil
	//	}
	//	// the same name chart exists, but allowed to override
	//	glog.V(2).Infof("the same name chart %s:%s exists on the destination registry and the \"override\" is set to true, continue...",
	//		name, version)
	//}

	// copy the chart between the source and destination registries
	chart, err := t.downloadChart(name, version)
	if err != nil {
		glog.Errorf("failed to download the chart %s:%s: %v", name, version, err)
		return err
	}
	defer chart.Close()

	if err = t.uploadChart(name, version, chart); err != nil {
		glog.Errorf("failed to upload the chart %s:%s: %v", name, version, err)
		return err
	}

	glog.V(2).Infof("copy %s:%s(source registry) to %s:%s(destination registry) completed",
		name, version, name, version)

	return nil
}

func mustHarborClient(repoAddr string)(*harbor.Client, error) {
	client := harbor.NewClient(nil, repoAddr, "admin", "Harbor12345")
	opt := harbor.ListProjectsOptions{Name: utils.DefaultProjectName}
	_, _, errs := client.Projects.ListProject(&opt)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return client, nil
}

func NewTransfer(srcAddr, dstAddr string) (*Transfer, error) {
	transport := utils.GetHTTPTransport(true)
	transport.Proxy =  func(req *http.Request) (*url.URL, error) {
		req.SetBasicAuth("admin", "Harbor12345")
		return nil, nil
	}

	sc, err := mustHarborClient(srcAddr)
	if err != nil {
		return nil, err
	}
	dc, err := mustHarborClient(dstAddr)
	if err != nil {
		return nil, err
	}

	return &Transfer{
		srcAddr : srcAddr,
		dstAddr : dstAddr,
		SrcRepo : sc,
		DstRepo : dc,
		client: common_http.NewClient(
			&http.Client{
				Transport: transport,
			},),
	}, nil
}