package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/portbinding"
	"net/http"
	"log"
)

var bindingsList []*portbinding.PortBinding = []*portbinding.PortBinding{}

func PortBindingCreate(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")
	binding, err := portbinding.New(port, host, hostPort)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	bindingsList = append(bindingsList, binding)
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
	port := r.FormValue("port")
	i := -1
	for j, binding := range bindingsList {
		if binding.Port == port {
			err := binding.Stop()
			if err != nil {
				w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
				return
			}
			i = j
			log.Println("stopped " + binding.String())
		}
	}
	if i == -1 {
		w.Write([]byte("{}"))
		return
	}
	if i != -1 {
		binding := bindingsList[i]
		content, err := json.Marshal(binding)
		if err != nil {
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			return
		}
		w.Write(content)
		copy(bindingsList[i:], bindingsList[i+1:])
		bindingsList[len(bindingsList)-1] = nil
		bindingsList = bindingsList[:len(bindingsList)-1]
	}
}
