package main

import (
	"flag"
	"strings"
	"strconv"
	"github.com/mmpei/gossip/model"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/mmpei/gossip/handler"
	"time"
	"github.com/mmpei/gossip/sync"
)

type sliceFlagValue []string

func (s *sliceFlagValue) Set(val string) error {
	*s = sliceFlagValue(strings.Split(val, ","))
	return nil
}

func (s *sliceFlagValue) String() string {
	*s = sliceFlagValue(strings.Split("20202,20203", ","))
	return ""
}

func main() {
	var port int
	var seeds sliceFlagValue
	flag.Var(&seeds, "seeds", "the seeds of gossip")
	flag.IntVar(&port, "port", 0, "server port")
	flag.Parse()

	ss := []string(seeds)
	var seedPorts []int
	for _, seed := range ss {
		s, _ := strconv.Atoi(seed)
		seedPorts = append(seedPorts, s)
		// add seed
		handler.Manager.Set(
			&model.PeerInfo{
				PeerId: s,
				Version: 0,
				Seed: true,
				Alive: true,
				SyncTime: time.Now(),
			},
		)
	}

	model.Self = model.NewPeer(port)
	// add self
	handler.Manager.Set(model.Self)

	syncManager := sync.NewSyncManager()
	ticker:=time.NewTicker(time.Second*2)
	go func() {
		for range ticker.C {
			syncManager.Announce()
		}
	}()

	listenAddr := fmt.Sprintf("0.0.0.0:%d", port)
	router := mux.NewRouter()
	router.HandleFunc("/metrics", handler.Metrics).Methods("GET")
	router.HandleFunc("/announce", handler.Announce).Methods("POST")
	http.ListenAndServe(listenAddr, router)
}
