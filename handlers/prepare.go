package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soider/d"
	"github.com/soider/schnur/ssh"
	"github.com/soider/schnur/targets/manager"
	"net/http"
)

type PrepareHandler struct {
	tm *manager.TargetsManager
}

func NewPrepareHandler(tm *manager.TargetsManager) *PrepareHandler {
	return &PrepareHandler{tm: tm}
}

type response struct {
	VncTarget string   `json:"vnc"`
	Output    []string `json:"output"`
}

func (ch PrepareHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqData := mux.Vars(req)
	targetName := reqData["target"]
	target, err := ch.tm.Target(targetName)
	if err != nil {
		http.Error(rw, "Error: "+err.Error(), 500)
		return
	}
	d.D(target)
	sshCfg := target.SshConfig()
	executor := ssh.New(sshCfg)
	cmds, err := target.GetCmd()
	if err != nil {
		http.Error(rw, "Can't prepare cmd for "+target.Name+": "+err.Error(), 500)
		return
	}
	var result []string
	for _, cmd := range cmds {
		output, err := executor.RunRemoteCMD(cmd)
		if err != nil {
			http.Error(rw, "Can't execute cmd: "+cmd+": "+err.Error()+" on "+target.Name, 500)
			return
		}
		result = append(result, output)
	}
	rw.Header().Set("Content-Type", "application/json")
	resp := response{
		Output:    result,
		VncTarget: fmt.Sprintf("%s:%s", target.GetVncAddress(), target.VncPort),
	}
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(resp)
	if err != nil {
		http.Error(rw, "Error while encoding answer: "+err.Error(), 500)
	}

}
