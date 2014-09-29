package handler

import (
	"encoding/json"
	h "github.com/romanoff/vanguard/host"
	"github.com/romanoff/vanguard/portbinding"
	"log"
	"net/http"
)

func getPortBinding(port string, createNew bool) (int, *portbinding.PortBinding, error) {
	host, err := h.GetCurrentHost()
	if err != nil {
		return -1, nil, err
	}
	for i, b := range host.PortBindings {
		if b.Port == port {
			return i, b, nil
		}
	}
	if createNew {
		binding, err := portbinding.New(port)
		if err != nil {
			return -1, nil, err
		}
		host.PortBindings = append(host.PortBindings, binding)
		return len(host.PortBindings) - 1, binding, nil
	}
	return -1, nil, nil

}

func PortBindingCreate(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")
	currentHost, err := h.GetCurrentHost()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	_, binding, err := getPortBinding(port, true)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	binding.AddBackend(host, hostPort)
	err = currentHost.Persist()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	go binding.Start()
	content, err := json.Marshal(binding)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	log.Println("exposed " + binding.String())
	w.Write(content)
}

func PortBindingsIndex(w http.ResponseWriter, r *http.Request) {
	host, err := h.GetCurrentHost()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	content, err := json.Marshal(host.PortBindings)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write(content)
}

func PortBindingDelete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	port := params.Get(":port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")

	i, binding, err := getPortBinding(port, false)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}

	if binding == nil {
		w.Write([]byte("{}"))
		return
	}
	if host != "" && hostPort != "" {
		binding.RemoveBackend(host, hostPort)
	}
	currentHost, err := h.GetCurrentHost()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	if len(binding.Backends) == 0 || (host == "" && hostPort == "") {
		binding.Stop()
		copy(currentHost.PortBindings[i:], currentHost.PortBindings[i+1:])
		currentHost.PortBindings[len(currentHost.PortBindings)-1] = nil
		currentHost.PortBindings = currentHost.PortBindings[:len(currentHost.PortBindings)-1]
	}
	w.Write([]byte("{\"success\": true}"))
	currentHost.Persist()
}
