package handlers

import (
	"net/http"
	"github.com/soider/schnur/targets/manager"
	"encoding/json"
)

type ListHandler struct {
	tm *manager.TargetsManager
}

func NewListHandler(tm *manager.TargetsManager) *ListHandler {
	return &ListHandler{tm: tm}
}

func (ch ListHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	targets, err := ch.tm.Targets()
	if err != nil {
		http.Error(rw, "Can't load targets: " + err.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(targets)
	if err != nil {
		http.Error(rw, "Can't encode targets: " + err.Error(), 500)
	}
}

