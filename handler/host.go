package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/host"
	"net/http"
)

func HostsIndex(w http.ResponseWriter, r *http.Request) {
	hosts, err := host.GetHosts()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	hostsJson, err := json.Marshal(hosts)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return		
	}
	w.Write(hostsJson)
}
