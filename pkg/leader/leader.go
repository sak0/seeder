package leader

import (
	"fmt"
	"time"
	"os"
	"github.com/toolkits/net"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/golang/glog"
)

var (
	ConsulURI = os.Getenv("CONSUL_ADDR") + ":" + os.Getenv("CONSUL_PORT")
	LeaderKey = "seederKey"
)


type LeaderElection struct {
	Client	 	*consulapi.Client
	TTL      	time.Duration
	key      	string
	CallBack 	func(stop chan struct{})
	StopCh		chan struct{}
}

func GetMyIPAddr() string {
	ips, err := net.IntranetIP()
	if err != nil {
		glog.V(2).Infof("can't get ip list: %v", err.Error())
		return ""
	}
	return ips[0]
}


func (le *LeaderElection)watchKey(key string, ch chan int)error{
	params := make(map[string]interface{})
	params["type"] = "key"
	params["key"] = key
	params["stale"] = false

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(index uint64, result interface{}) {
		if kvpair, ok := result.(*consulapi.KVPair); ok {
			if kvpair.Session == "" {
				fmt.Printf("The key %s's lock session is released.\n", key)
				ch <- 1
			}
		}
	}

	err = plan.Run(ConsulURI)
	if err != nil {
		return err
	}

	return nil
}

func (le *LeaderElection)Run(){
	var IdentityName = GetMyIPAddr()

	se := &consulapi.SessionEntry{
		Name:      IdentityName,
		TTL:       le.TTL.String(),
		LockDelay: time.Nanosecond,
	}
	for {
		sessionId, _, err :=  le.Client.Session().CreateNoChecks(se, nil)
		if err != nil {
			fmt.Printf("[%s] Create Session failed: %v\n", IdentityName, err)
			time.Sleep(10 * time.Second)
			continue
		}

		kvpair := &consulapi.KVPair{
			Key 	: LeaderKey,
			Value 	: []byte(IdentityName),
			Session : sessionId,
		}
		locked, _, err := le.Client.KV().Acquire(kvpair, nil)
		if err != nil {
			fmt.Printf("[%s] Acquire Lock key %s failed: %v", IdentityName, LeaderKey, err)
			continue
		}

		if !locked {
			fmt.Printf("[%s] is follower. Begin watch the key.\n", IdentityName)
			ch := make(chan int, 1)
			go le.watchKey(LeaderKey, ch)
			<-ch
			fmt.Printf("[%s] The lock is release, Again election.\n", IdentityName)
		} else {
			fmt.Printf("[%s] is leader now.\n", IdentityName)
			go le.CallBack(le.StopCh)
			le.Client.Session().RenewPeriodic(se.TTL, sessionId, nil, nil)
		}
	}
}
