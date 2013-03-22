package server

import (
	"github.com/ha/doozer"
	"github.com/ActiveState/doozerconfig"
	"github.com/ActiveState/log"
	"strings"
)

type clusterConfig struct {
	Endpoint string `doozer:"/cluster/config/endpoint"`
	CoreIP   string `doozer:"/cluster/config/mbusip"`
	NatsUri  string `doozer:"/proc/cloud_controller/config/mbus"`
}

var Config *clusterConfig

// IsMicro returns true if the cluster is configured as a micro cloud.
func (c *clusterConfig) IsMicro() bool {
	return strings.Contains(c.NatsUri, "/127.0.0.1:")
}

func Init(conn *doozer.Conn, rev int64) {
	Config = new(clusterConfig)
	cfg := doozerconfig.New(conn, rev, Config, "")
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}

	go cfg.Monitor("/cluster/config/*", func(change *doozerconfig.Change, err error) {
		if err != nil {
			log.Errorf("Unable to process cluster config change in doozer: %s", err)
			return
		}
	})
}
