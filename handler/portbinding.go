package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/portbinding"
	"log"
	"net/http"
)

var bindingsList []*portbinding.PortBinding = []*portbinding.PortBinding{}

func getPortBinding(port string, createNew bool) (int, *portbinding.PortBinding, error) {
	for i, b := range bindingsList {
		if b.Port == port {
			return i, b, nil
		}
	}
	if createNew {
		binding, err := portbinding.New(port)
		if err != nil {
			return -1, nil, err
		}
		bindingsList = append(bindingsList, binding)
		return len(bindingsList) - 1, binding, nil
	}
	return -1, nil, nil

}

func PortBindingCreate(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")
	_, binding, err := getPortBinding(port, true)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	binding.AddBackend(host, hostPort)
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
	content, err := json.Marshal(bindingsList)
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
	if len(binding.Backends) == 0 || (host == "" && hostPort == "") {
		binding.Stop()
		copy(bindingsList[i:], bindingsList[i+1:])
		bindingsList[len(bindingsList)-1] = nil
		bindingsList = bindingsList[:len(bindingsList)-1]
	}
	w.Write([]byte("{\"success\": true}"))
}
