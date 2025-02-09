package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/toolkits/net"
	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
	"runtime"
	"runtime/pprof"
	"net/http"
	"crypto/tls"
)

const (
	DeRegisterInterval 	= "2m"
	HealthCheckTimeout	= "10s"
	HealthCheckInterval	= "15s"

	DefaultProjectName 	= "edge-cloud"
)

var (
	defaultHTTPTransport, secureHTTPTransport, insecureHTTPTransport *http.Transport
	MyNodeName 	string
	HarborUser 	string
	HarborPass 	string
	MyRole 		string
)

func SetNodeName(myName, myRole string) {
	MyNodeName = myName
	MyRole = myRole
}

func GetMyNodeName() string {
	return MyNodeName
}

func init() {
	defaultHTTPTransport = &http.Transport{}

	secureHTTPTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	insecureHTTPTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func GetHTTPTransport(insecure ...bool) *http.Transport {
	if len(insecure) == 0 {
		return defaultHTTPTransport
	}
	if insecure[0] {
		return insecureHTTPTransport
	}
	return secureHTTPTransport
}

func ArrayIn(item string, arr []string) bool {
	for _, a := range arr {
		if item == a {
			return true
		}
	}
	return false
}

func GetMyIpAddr() (string, error) {
	ips, err := net.IntranetIP()
	if err != nil {
		return "", err
	}
	return ips[0], nil
}

func MustGetMyIpAddr() string {
	ips, err := net.IntranetIP()
	if err != nil {
		panic(err)
	}
	return ips[0]
}

func HarborAuth() error {
	harborUser := os.Getenv("HARBOR_USER")
	harborPasswd := os.Getenv("HARBOR_PASSWD")
	if harborUser == "" || harborPasswd == "" {
		return fmt.Errorf("Env HARBOR_USER and HARBOR_PASSWD should exists.\n")
	}

	HarborUser = harborUser
	HarborPass = harborPasswd

	return nil
}

func ServiceRegister(myName string, myPort int, healthURL string) error {
	consulAddr := os.Getenv("CONSUL_ADDR")
	consulPort := os.Getenv("CONSUL_PORT")
	if consulAddr == "" || consulPort == "" {
		return fmt.Errorf("Env CONSUL_ADDR and CONSUL_PORT should exists.\n")
	}

	consulConfig := &api.Config{
		Address: consulAddr + ":" + consulPort,
	}
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}

	myIp, err := GetMyIpAddr()
	if err != nil {
		return err
	}

	glog.V(2).Infof("Register %s:%s (health:%s)to consul %v",
		myIp, strconv.Itoa(myPort), healthURL, consulConfig.Address)

	myCheck := api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: DeRegisterInterval,
		Timeout: HealthCheckTimeout,
		Interval: HealthCheckInterval,
		HTTP: fmt.Sprintf("http://%s:%s/%s", myIp, strconv.Itoa(myPort), healthURL),
	}

	register := api.AgentServiceRegistration{
		ID: fmt.Sprintf("%s_%s_%s", myName, myIp, strconv.Itoa(myPort)),
		Name: myName,
		Port: myPort,
		Address: myIp,
		Check: &myCheck,
	}
	return client.Agent().ServiceRegister(&register)
}

func ConsumeMem() uint64 {
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Sys
}

func DoResourceMonitor() {
	m := pprof.Lookup("goroutine")
	memStats := ConsumeMem()
	glog.V(3).Infof("Resource monitor: [%d goroutines] [%.3f kb]", m.Count(), float64(memStats)/1e3)
}