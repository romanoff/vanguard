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

func HostShow(w http.ResponseWriter, r *http.Request) {
	h, err := host.GetCurrentHost()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	hostJson, err := json.Marshal(h)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write(hostJson)
}
