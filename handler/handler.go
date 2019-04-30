package handler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/mmpei/gossip/model"
)

func Metrics(w http.ResponseWriter, r *http.Request) {
	list := Manager.Get()

	list = append(list, *model.Self)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		fmt.Errorf("json encode response: %s", err)
	}
}

func Announce(w http.ResponseWriter, r *http.Request) {
	req := []model.PeerInfo{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("json decode request: %s \n", err)
		w.WriteHeader(500)
		return
	}

	for index := range req {
		Manager.Set(&req[index])
	}
	w.Write([]byte("ok"))
}
