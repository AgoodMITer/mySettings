package sync

import (
	"net/http"
	"time"
	"github.com/mmpei/gossip/handler"
	"github.com/mmpei/gossip/model"
	"math/rand"
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

const SyncPeers = 3

type SyncManager struct {
	client http.Client
}

func NewSyncManager() *SyncManager {
	return &SyncManager{
		client: http.Client{
			Timeout: 10*time.Second,
		},
	}
}

func (sm *SyncManager) Announce() {
	model.Self.Version++
	model.Self.SyncTime = time.Now()
	handler.Manager.Set(model.Self)
	peers := handler.Manager.Get()
	eps, _ := pickPeers(peers, SyncPeers)

	peers = append(peers, *model.Self)
	bytesData, err := json.Marshal(peers)
	if err != nil {
		fmt.Println(err.Error() )
		return
	}

	fmt.Printf("announce to remote: version=%d, syncPeers=%d, peers=%d \n", model.Self.Version, len(eps), len(peers))
	for _, ep := range eps {
		url := fmt.Sprintf("http://localhost:%d/announce", ep.PeerId)
		reader := bytes.NewReader(bytesData)
		id := ep.PeerId
		go func () {
			request, err := http.NewRequest("POST", url, reader)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
			resp, err := sm.client.Do(request)
			if err != nil { // set remote peer down
				fmt.Println(err.Error())
				handler.Manager.SetDown(id)
				return
			}
			ioutil.ReadAll(resp.Body)
		}()
	}
}

// pickPeers pick count peers randomly from peers
func pickPeers(peers []model.PeerInfo, count int) ([]model.PeerInfo, error) {
	if len(peers) <= count {
		return peers, nil
	}

	sum := len(peers)
	eps := []model.PeerInfo{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	indexs := []int{}

PickLoop:
	for i:=0; i < count; {
		index := r.Intn(sum)
		for _, ind := range indexs {
			if ind == index {
				continue PickLoop
			}
		}
		indexs = append(indexs, index)
		eps = append(eps, peers[index])
		i++
	}
	return eps, nil
}
